package logic

import (
	"context"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/repo"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

var (
	FRIEND_NOT_EXIST    = response.MsgCode{51002, "好友不存在"}
	JUDGE_FRIEND        = response.MsgCode{51003, "验证好友关系失败"}
	GET_HISTORY_MESSAGE = response.MsgCode{51001, "获取聊天记录失败"}
	SAVE_MESSAGE_FAILED = response.MsgCode{51004, "存储消息失败"}
	PUSH_MESSAGE_FAILED = response.MsgCode{51005, "推送消息失败"}

	mutex   sync.Mutex
	clients = make(map[int64]*websocket.Conn)
)

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

// SendMessage
//
//	@Description: 登录用户发送信息给指定好友
//	@receiver l
//	@param ctx
//	@param SenderID
//	@param msg
//	@return err
func (l *Chatlogic) SendMessage(ctx context.Context, SenderID int64, msg types.WSMessage) (err error) {
	defer utils.RecordTime(time.Now())()
	//验证好友关系
	exist, err := repo.NewFriendRepo(global.DB).JudgeFriend(SenderID, msg.To)
	if err != nil {
		zlog.CtxErrorf(ctx, "验证好友关系失败: %v", err)
		return response.ErrResp(err, JUDGE_FRIEND)
	}
	if exist == false {
		zlog.CtxErrorf(ctx, "用户%d与用户%d不是好友关系", SenderID, msg.To)
		return response.ErrResp(nil, FRIEND_NOT_EXIST)
	}

	//开启事务，确保存储消息和推送消息同时成功或失败
	tx := global.DB.Begin()

	// 存储消息
	err = repo.NewChatRepo(global.DB).SaveMessage(SenderID, msg)
	if err != nil {
		zlog.CtxErrorf(ctx, "存储消息失败: %v", err)
		tx.Rollback()
		return response.ErrResp(err, SAVE_MESSAGE_FAILED)
	}

	// 推送消息
	err = PushMessage(msg.To, types.WSMessageResp{
		From:    SenderID,
		Content: msg.Content,
		Time:    time.Now().Unix(), // 当前时间的时间戳
	})
	if err != nil {
		zlog.CtxErrorf(ctx, "推送消息失败: %v", err)
		tx.Rollback()
		return response.ErrResp(err, PUSH_MESSAGE_FAILED)
	}
	tx.Commit()

	return nil
}

// PushMessage
//
//	@Description: 把消息发送给接收者
//	@param receiverID
//	@param resp
func PushMessage(receiverID int64, resp types.WSMessageResp) error {
	mutex.Lock()
	defer mutex.Unlock()

	if conn, exists := clients[receiverID]; exists {
		err := conn.WriteJSON(resp)
		if err != nil {
			zlog.Errorf("PushMessage error: %v", err)
			return err
		}
	}
	return nil
}
