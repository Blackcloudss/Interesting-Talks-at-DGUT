package logic

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/connect"

	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
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
)

const (
	RETRYCOUNT = "retry_count"
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

// 消息队列相关配置
const (
	MaxBatchSize  = 50  // 批量消息最大条数，平衡吞吐量与延迟
	FlushInterval = 100 // 批量处理间隔(ms)，控制消息实时性
)

var (
	MsgQueue    = make(chan types.WSMessageResp, 1000) // 缓冲消息队列，应对突发流量
	buffer      []types.WSMessageResp
	ackChannels = sync.Map{} // 消息确认通道映射，实现消息状态跟踪
)

// SendMessage 消息发送入口
// 解决：消息可靠性、好友验证、存储与推送解耦
func (l *Chatlogic) SendMessage(ctx context.Context, SenderID int64, WSMsg types.WSMessage, CM *connect.ConnectionManager) (err error) {
	defer utils.RecordTime(time.Now())()
	//	好友关系验证
	exist, err := repo.NewFriendRepo(global.DB).JudgeFriend(SenderID, WSMsg.To)
	if err != nil {
		zlog.CtxErrorf(ctx, "验证好友关系失败: %v", err)
		return response.ErrResp(err, JUDGE_FRIEND)
	}
	if exist == false {
		zlog.CtxErrorf(ctx, "用户%d与用户%d不是好友关系", SenderID, WSMsg.To)
		return response.ErrResp(nil, FRIEND_NOT_EXIST)
	}

	// 存储消息时初始化状态
	status := global.DELIVERED
	if _, exists := CM.Clients[WSMsg.To]; !exists {
		status = global.OFFLINE
	}

	// 消息存储
	err = repo.NewChatRepo(global.DB).SaveMessage(SenderID, WSMsg, status)
	if err != nil {
		zlog.CtxErrorf(ctx, "存储消息失败: %v", err)
		return response.ErrResp(err, SAVE_MESSAGE_FAILED)
	}

	// 构建响应消息
	respMsg := types.WSMessageResp{
		From:    SenderID,
		To:      WSMsg.To,
		Content: WSMsg.Content,
		Time:    time.Now().Unix(),
		MsgID:   WSMsg.MsgID,
		Type:    global.MESSAGE,
	}

	// 根据消息状态进行推送
	if status == global.DELIVERED {
		MsgQueue <- respMsg
	} else {
		if err = repo.NewChatRepo(global.DB).SaveOfflineMessage(respMsg); err != nil {
			zlog.CtxErrorf(ctx, "保存离线消息失败: %v", err)
		}
	}

	// 批处理模式消息队列
	var initOnce sync.Once
	initOnce.Do(func() {
		go ProcessMsgQueue(CM)
	})

	// 在flushMessages中增加确认等待
	go func(msgID string) {
		select {
		case <-time.After(30 * time.Second):
			if !checkAck(msgID) {
				MsgQueue <- types.WSMessageResp{
					From:    SenderID,
					To:      WSMsg.To,
					Content: WSMsg.Content,
					Time:    time.Now().Unix(),
					MsgID:   msgID,
				}
			}
		case <-ackChan(msgID):
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

	CM.Mutex.Lock()
	defer CM.Mutex.Unlock()

	for _, msg := range msgs {
		if conn, exists := CM.Clients[msg.To]; exists {
			if err = conn.WriteMessage(websocket.BinaryMessage, buf.Bytes()); err == nil {
				// 记录消息发送时间
				ackChannels.Store(msg.MsgID, time.Now())
			}
		} else {
			err = repo.NewChatRepo(global.DB).SaveOfflineMessage(msg)
			if err != nil {
				zlog.Errorf("保存离线消息失败: %v", err)
			}
		}
	}
	return nil
}

func ackChan(msgID string) <-chan struct{} {
	ch, _ := ackChannels.LoadOrStore(msgID, make(chan struct{}))
	return ch.(chan struct{})
}

func checkAck(msgID string) bool {
	_, ok := ackChannels.Load(msgID)
	return ok
}

func updateMsgStatus(msgID, status string) {
	if err := repo.NewChatRepo(global.DB).UpdateMessageStatus(msgID, status); err != nil {
		zlog.Errorf("更新消息状态失败: %v", err)
	}

	// 自动清理超过3次重试的消息
	if status == global.DELIVERED {
		var retryCount int
		repo.NewChatRepo(global.DB).DB.Model(&model.Message{}).
			Where(fmt.Sprintf("%v = ?", global.MSGID), msgID).
			Pluck(fmt.Sprintf("%v = ?", RETRYCOUNT), &retryCount)

		if retryCount >= 3 {
			repo.NewChatRepo(global.DB).DB.
				Where(fmt.Sprintf("%v = ?", global.MSGID), msgID).
				Delete(&model.Message{})
		}
	}

	ackChannels.Delete(msgID)
}

func (l *Chatlogic) GetOfflineMessages(userID int64) ([]types.WSMessageResp, error) {
	return repo.NewChatRepo(global.DB).GetPendingMessages(userID)
}
