package types

import "time"

// 通过Websocket获取前端发过来的消息
type WSMessage struct {
	To      int64  `json:"to" binding:"required"` //目标好友ID -- 接受者ID
	Content string `json:"content"`               //要发送给好友的内容
}

// 将Websocket错误信息返回给前端
type WSError struct {
	ERROR string `json:"error"`
}

// 将Websocket消息发给接收者
type WSMessageResp struct {
	From    int64  `json:"from"`    //发送者ID
	Content string `json:"content"` //发送内容
	Time    int64  `json:"time"`    //发送时间 -- 时间戳
}

// 获取历史聊天消息 入参
type GetMessageReq struct {
	ReceiverID int64 `json:"receiver_id" binding:"required"` //接收者ID
	Page       int   `json:"page"`                           //页码
	Size       int   `json:"size"`                           //每页大小
}

// 历史聊天消息
type MessageHistory struct {
	Sender   int64     `json:"sender"`   //发送者ID
	Receiver int64     `json:"receiver"` //接收者ID
	Content  string    `json:"content"`  //发送内容
	Time     time.Time `json:"time"`     //发送时间
}

// 获取历史聊天消息 出参
type GetMessageResp struct {
	MessagesHistory []MessageHistory `json:"messages"` //历史聊天消息
}
