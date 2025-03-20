package types

import "time"

/*
数据结构定义模块
主要功能：
1. 定义API请求/响应格式
2. 定义持久化数据结构
3. 定义WebSocket协议格式
*/

// WSMessage 客户端发送消息结构
// 包含：接收方ID、内容、唯一消息ID、消息类型
type WSMessage struct {
	To      int64  `json:"to" binding:"required"` // 目标好友ID
	Content string `json:"content"`               // 消息内容
	MsgID   string `json:"msg_id"`                // 消息唯一ID
	Type    string `json:"type"`                  // 消息类型（message/ack）
}

// 将Websocket错误信息返回给前端
type WSError struct {
	ERROR string `json:"error"`
}

// WSMessageResp 服务端推送结构
// 包含：发送方信息、时间戳、状态标识
type WSMessageResp struct {
	From    int64  `json:"from"`    // 发送者ID
	To      int64  `json:"to"`      // 接收者ID
	Content string `json:"content"` // 消息内容
	Time    int64  `json:"time"`    // 时间戳
	MsgID   string `json:"msg_id"`  // 消息唯一ID
	Type    string `json:"type"`    // 消息类型（message/ack）
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
