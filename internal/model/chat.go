package model

// @Title        friend.go
// @Description
// @Create       XdpCs 2025-03-20 上午12:14
// @Update       XdpCs 2025-03-20 上午12:14

// 聊天记录
type Message struct {
	CommonModel
	Sender   int64  `gorm:"column:sender;type:bigint;comment:'发送者ID'"`
	Receiver int64  `gorm:"column:receiver;type:bigint;comment:'接受者ID'"`
	Content  string `gorm:"column:content;type:text;comment:'聊天内容'" json:"content"`
}
