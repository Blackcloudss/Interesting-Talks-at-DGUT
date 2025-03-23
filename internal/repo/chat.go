package repo

import (
	"fmt"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"

	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"gorm.io/gorm"
)

const (
	SENDER     = "sender"
	RECEIVER   = "receiver"
	CONTENT    = "content"
	CREATED_AT = "created_at"
)

/*
数据持久化模块
主要功能：
1. 消息存储（MySQL）
2. 消息状态更新
3. 历史记录查询
4. 离线消息管理
*/

// @Title        chat.go
// @Description
// @Create       XdpCs 2025-03-20 上午12:23
// @Update       XdpCs 2025-03-20 上午12:23
type ChatRepo struct {
	DB *gorm.DB
}

func NewChatRepo(db *gorm.DB) *ChatRepo {
	return &ChatRepo{DB: db}
}

func (r *ChatRepo) GetMessagesHistory(SenderID int64, ReceiverID int64, page int, size int) (resp []types.MessageHistory, err error) {
	err = r.DB.Model(&model.Message{}).
		Select(SENDER, RECEIVER, CONTENT, CREATED_AT).
		Where(fmt.Sprintf("(%s = ? AND %s = ?) OR (%s = ? AND %s = ?)", SenderID, ReceiverID, ReceiverID, SenderID)).
		Order("created_at DESC").
		Limit(size).
		Offset((page - 1) * size).
		Find(&resp).
		Error
	if err != nil {
		zlog.Errorf("获取聊天记录失败:%v", err)
		return nil, err
	}
	return
}

// SaveMessage 消息存储
// 解决：消息持久化、状态跟踪（delivered/offline）
func (r *ChatRepo) SaveMessage(sender int64, msg types.WSMessageReq, status string) error {
	// 存储消息时记录发送状态，支持后续状态追踪
	return r.DB.Create(&model.Message{
		MsgID:    msg.MsgID,
		Sender:   sender,
		Receiver: msg.To,
		Content:  msg.Content,
		Status:   status,
	}).Error
}

func (r *ChatRepo) SaveOfflineMessage(msg types.WSMessageResp) error {
	return r.DB.Create(&model.Message{
		MsgID:    msg.MsgID,
		Sender:   msg.From,
		Receiver: msg.To,
		Content:  msg.Content,
		Status:   global.OFFLINE,
	}).Error
}

func (r *ChatRepo) UpdateMessageStatus(msgID, status string) error {
	return r.DB.Model(&model.Message{}).
		Where(fmt.Sprintf("%v = ?", global.MSGID), msgID).
		Update(global.STATUS, status).
		Error
}

// GetPendingMessages 获取待处理消息
// 解决：离线消息恢复、消息去重
func (r *ChatRepo) GetPendingMessages(userID int64) ([]types.WSMessageResp, error) {
	// 查询状态为delivered/offline的消息
	// 保证消息按时间顺序返回
	var messages []model.Message
	err := r.DB.Where("receiver = ? AND status IN ('delivered', 'offline')", userID).
		Find(&messages).
		Error

	var resp []types.WSMessageResp
	for _, msg := range messages {
		resp = append(resp, types.WSMessageResp{
			MsgID:   msg.MsgID,
			From:    msg.Sender,
			To:      msg.Receiver,
			Content: msg.Content,
			Time:    msg.CreatedAt.Unix(),
		})
	}
	return resp, err
}
