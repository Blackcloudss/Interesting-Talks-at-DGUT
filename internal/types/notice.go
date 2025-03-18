package types

import "github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"

// CreateNoticeReq 创建通知请求结构体
type CreateNoticeReq struct {
	Title   string `json:"title"`   // 通知标题
	Content string `json:"content"` // 通知内容
	UserID  int64  `json:"user_id"` // 用户ID
}

// CreateNoticeResp 创建通知响应结构体
type CreateNoticeResp struct {
	Notice model.Notice `json:"notice"`
}

// UpdateNoticeReq 更新通知请求结构体
type UpdateNoticeReq struct {
	ID      int64  `json:"id"`      // 通知ID
	UserID  int64  `json:"user_id"` // 用户ID
	Title   string `json:"title"`   // 通知标题
	Content string `json:"content"` // 通知内容
}

// UpdateNoticeResp 更新通知响应结构体
type UpdateNoticeResp struct {
	Notice model.Notice `json:"notice"`
}

// DeleteNoticeReq 删除通知请求结构体
type DeleteNoticeReq struct {
	ID     int64 `json:"id"`      // 通知ID
	UserID int64 `json:"user_id"` // 用户ID
}

// DeleteNoticeResp 删除通知响应结构体
type DeleteNoticeResp struct {
	Success bool `json:"success"`
}

// GetNoticeReq 获取通知请求结构体
type GetNoticeReq struct {
	ID int64 `json:"id"` // 通知ID
}

// GetNoticeResp 获取通知响应结构体
type GetNoticeResp struct {
	Notice model.Notice `json:"notice"`
}
