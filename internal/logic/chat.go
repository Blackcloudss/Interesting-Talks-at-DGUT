package logic

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/connect"

	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/repo"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils"
	"github.com/gorilla/websocket"
)

var (
	FRIEND_NOT_EXIST    = response.MsgCode{51002, "好友不存在"}
	JUDGE_FRIEND        = response.MsgCode{51003, "验证好友关系失败"}
	GET_HISTORY_MESSAGE = response.MsgCode{51001, "获取聊天记录失败"}
	SAVE_MESSAGE_FAILED = response.MsgCode{51004, "存储消息失败"}

	MsgQueue    = make(chan types.WSMessageResp, 1000) // 缓冲消息队列，应对突发流量
	buffer      []types.WSMessageResp                  // 消息缓冲区，用于存储待处理的消息
	ackChannels = sync.Map{}                           // "消息确认"通道映射，实现消息状态跟踪
)

const (
	//消息队列相关配置
	MaxBatchSize  = 50  // 批量消息最大条数，平衡吞吐量与延迟
	FlushInterval = 100 // 批量处理间隔(ms)，控制消息实时性
)

/*
聊天核心业务逻辑模块
主要功能：
1. 消息发送验证（好友关系、内容合法性）
2. 消息存储与状态管理（在线/离线）
3. 消息批量压缩传输
4. 消息确认与重试机制
5. 离线消息管理
*/

// @Title        friend.go
// @Description
// @Create       XdpCs 2025-03-19 下午3:07
// @Update       XdpCs 2025-03-19 下午3:07
type Chatlogic struct {
}

func NewChatlogic() *Chatlogic {
	return &Chatlogic{}
}

// GetMessagesHistory 获取聊天记录
func (l *Chatlogic) GetMessagesHistory(ctx context.Context, SenderID int64, req types.GetMessageReq) (resp types.GetMessageResp, err error) {
	defer utils.RecordTime(time.Now())()
	// 如果用户没有页码和页数，则默认为第一页和每页10条数据
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 10
	}
	// 获取聊天记录
	resp.MessagesHistory, err = repo.NewChatRepo(global.DB).GetMessagesHistory(SenderID, req.ReceiverID, req.Page, req.Size)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取聊天记录失败: %v", err)
		return resp, response.ErrResp(err, GET_HISTORY_MESSAGE)
	}
	return resp, nil
}

// SendMessage 消息发送入口
// 解决：消息可靠性、好友验证、存储与推送解耦
func (l *Chatlogic) SendMessage(ctx context.Context, SenderID int64, WSMsg types.WSMessageReq, CM *connect.ConnectionManager) (err error) {
	defer utils.RecordTime(time.Now())()
	//	好友关系验证
	exist, err := repo.NewFriendRepo(global.DB).JudgeFriend(SenderID, WSMsg.To)
	if err != nil {
		zlog.CtxErrorf(ctx, "验证好友关系失败: %v", err)
		return response.ErrResp(err, JUDGE_FRIEND)
	}
	// 如果不是好友关系，则返回错误
	if exist == false {
		zlog.CtxErrorf(ctx, "用户%d与用户%d不是好友关系", SenderID, WSMsg.To)
		return response.ErrResp(nil, FRIEND_NOT_EXIST)
	}

	// 存储消息时初始化状态 DELIVERED 在线
	status := global.DELIVERED
	// 如果好友不在线，则消息状态为 OFFLINE 离线
	if _, exists := CM.Clients[WSMsg.To]; !exists {
		status = global.OFFLINE
	}

	// 构建响应消息 -- 普通消息
	respMsg := types.WSMessageResp{
		From:    SenderID,
		To:      WSMsg.To,
		Content: WSMsg.Content,
		Time:    time.Now().Unix(),
		MsgID:   WSMsg.MsgID,
		Type:    global.MESSAGE,
	}

	// 消息存储
	err = repo.NewChatRepo(global.DB).SaveMessage(respMsg, status)
	if err != nil {
		zlog.CtxErrorf(ctx, "存储消息失败: %v", err)
		return response.ErrResp(err, SAVE_MESSAGE_FAILED)
	}

	// 根据消息状态进行推送
	// 如果对方在线，则直接推送到 缓冲消息队列
	if status == global.DELIVERED {
		MsgQueue <- respMsg
	}

	// 批处理模式消息队列
	var initOnce sync.Once
	initOnce.Do(func() {
		go ProcessMsgQueue(CM)
	})

	// 在flushMessages中增加确认等待
	// 通过ACK机制确保消息必达（若30秒未收到ACK会触发重发）
	go func(msgID string) {
		select {
		case <-time.After(30 * time.Second):
			if !checkAck(msgID) { // 未收到ACK，重新入队
				// 获取消息重试信息
				message, err := repo.NewChatRepo(global.DB).GetMessageRetryInfo(msgID)
				if err != nil {
					zlog.Errorf("获取消息重试信息失败: %v", err)
					return
				}

				// 检查重试次数是否超过限制
				if message.RetryCount >= message.MaxRetries {
					zlog.Warnf("消息%s达到最大重试次数(%d)，将被删除", msgID, message.RetryCount)
					if err := repo.NewChatRepo(global.DB).DeleteMessage(msgID); err != nil {
						zlog.Errorf("删除消息失败: %v", err)
					}
					return
				}

				// 增加重试计数
				if err := repo.NewChatRepo(global.DB).IncrementRetryCount(msgID); err != nil {
					zlog.Errorf("重试计数更新失败: %v", err)
					return
				}

				// 重新入队
				MsgQueue <- types.WSMessageResp{
					From:       message.Sender,
					To:         message.Receiver,
					Content:    message.Content,
					Time:       time.Now().Unix(),
					MsgID:      msgID,
					Type:       global.MESSAGE,
					RetryCount: message.RetryCount + 1,
					MaxRetries: message.MaxRetries,
				}
			}
		case <-ackChan(msgID): // <-ch  // 会永久阻塞，直到通道被关闭或发送数据（此处通道类型是 struct{}，无数据发送）
			// 收到ACK， 更新消息状态为ACKNOWLEDGED
			updateMsgStatus(msgID, global.ACKNOWLEDGED)
		}
	}(WSMsg.MsgID)

	return nil
}

// ProcessMsgQueue 消息队列处理器
// 解决：消息批量处理、流量整形、压缩优化
func ProcessMsgQueue(CM *connect.ConnectionManager) {
	//创建一个定时器，用于每100毫秒执行一次批处理操作
	flushTimer := time.NewTicker(FlushInterval * time.Millisecond)
	defer flushTimer.Stop()
	for {
		select {
		// 从消息队列中获取消息
		case msg := <-MsgQueue:
			// 将消息添加到缓冲区
			buffer = append(buffer, msg)
			// 如果缓冲区达到最大批处理数量，则执行批处理操作
			if len(buffer) >= MaxBatchSize {
				err := flushMessages(buffer, CM)
				if err != nil {
					zlog.Errorf("MaxBatchSize_flushMessages error: %v", err)
				}
				buffer = nil
			}
			// 如果缓冲区时间到，则执行批处理操作
		case <-flushTimer.C:
			// 如果缓冲区不为空，则执行批处理操作
			if len(buffer) > 0 {
				err := flushMessages(buffer, CM)
				if err != nil {
					zlog.Errorf("flushTimer.C_flushMessages error: %v", err)
				}
				buffer = nil
			}
		}
	}
}

// flushMessages 消息批量发送
// 解决：网络效率优化（压缩）、接收方状态判断
func flushMessages(msgs []types.WSMessageResp, CM *connect.ConnectionManager) (err error) {
	// 使用 gzip 压缩批处理消息
	var buf bytes.Buffer
	// 创建 gzip 压缩器
	gz := gzip.NewWriter(&buf)
	// 设置压缩级别
	gz, err = gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	if err != nil {
		zlog.Errorf("创建gzip压缩器失败: %v", err)
		return err
	}
	// 将批处理消息编码为 JSON
	err = json.NewEncoder(gz).Encode(msgs)
	if err != nil {
		zlog.Errorf("批处理消息编码为 JSON 失败: %v", err)
		return
	}
	// 关闭 gzip 压缩器
	err = gz.Close()
	if err != nil {
		zlog.Errorf("关闭 gzip 压缩器失败: %v", err)
		return
	}
	// 广播消息
	CM.Mutex.Lock()
	defer CM.Mutex.Unlock()
	// 遍历所有在线用户
	for _, msg := range msgs {
		// 如果接收方在线，则发送消息
		if conn, exists := CM.Clients[msg.To]; exists {
			// 发送消息
			if err = conn.WriteMessage(websocket.BinaryMessage, buf.Bytes()); err == nil {
				// 记录消息发送时间
				ackChannels.Store(msg.MsgID, time.Now())
			}
		} else {
			// 如果接收方不在线，则存储离线消息
			err = repo.NewChatRepo(global.DB).SaveMessage(msg, global.OFFLINE)
			if err != nil {
				zlog.Errorf("保存离线消息失败: %v", err)
			}
		}
	}
	return nil
}

// ackChan 获取"消息确认"通道
func ackChan(msgID string) <-chan struct{} {
	//检查键是否存在，存在则返回值，否则存储新值。这里的值是一个通道。
	//始终返回通道对象（不论是否关闭）
	ch, _ := ackChannels.LoadOrStore(msgID, make(chan struct{}))
	return ch.(chan struct{})
}

// checkAck 检查消息是否已确认
func checkAck(msgID string) bool {
	_, ok := ackChannels.Load(msgID)
	return ok
}

// updateMsgStatus 更新消息状态
func updateMsgStatus(msgID, status string) {
	if err := repo.NewChatRepo(global.DB).UpdateMessageStatus(msgID, status); err != nil {
		zlog.Errorf("更新消息状态失败: %v", err)
		return
	}
	// 清理确认通道
	ackChannels.Delete(msgID)
}

// GetOfflineMessages 获取离线消息
func (l *Chatlogic) GetOfflineMessages(userID int64) (resp []types.WSMessageResp, err error) {
	resp, err = repo.NewChatRepo(global.DB).GetPendingMessages(userID)
	if err != nil {
		zlog.Errorf("获取离线消息失败: %v", err)
		return
	}
	return
}
