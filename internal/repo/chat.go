package repo

import (
	"fmt"
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
		Where(fmt.Sprintf("%s = ? AND %s = ?"), SenderID, ReceiverID).Or(fmt.Sprintf("%s = ? AND %s = ?"), ReceiverID, SenderID).
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

// SaveMessage
//
//	@Description: 保存聊天消息
//	@receiver r
//	@param SenderID
//	@param msg
//	@return err
func (r *ChatRepo) SaveMessage(SenderID int64, msg types.WSMessage) (err error) {
	// 保存消息
	err = r.DB.Create(&model.Message{
		Sender:   SenderID,
		Receiver: msg.To,
		Content:  msg.Content,
	}).Error
	if err != nil {
		zlog.Errorf("保存聊天消息失败:%v", err)
		return err
	}
	return nil
}
