package model

// @Title        friend.go
// @Description
// @Create       XdpCs 2025-03-20 上午12:14
// @Update       XdpCs 2025-03-20 上午12:14

/*
数据库模型定义模块
主要功能：
1. 定义ORM映射结构
2. 定义表结构约束
3. 定义业务状态枚举
*/

// Message 聊天消息模型
// 解决：数据持久化格式、查询优化
type Message struct {
	CommonModel
	MsgID      string `gorm:"column:msg_id;type:varchar(255);uniqueIndex;comment:'消息唯一ID'"`
	Sender     int64  `gorm:"column:sender;type:bigint;comment:'发送者ID'"`
	Receiver   int64  `gorm:"column:receiver;type:bigint;comment:'接受者ID'"`
	Content    string `gorm:"column:content;type:text;comment:'聊天内容'"`
	Status     string `gorm:"column:status;type:varchar(20);default:'delivered';comment:'消息状态'"`
	RetryCount int    `gorm:"column:retry_count;type:int;default:0;comment:'重试次数'"`
}
