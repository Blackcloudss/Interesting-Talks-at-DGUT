package types

// CreateNoticeReq 创建公告请求结构体
type CreateNoticeReq struct {
	Content string `json:"content"` // 公告内容
}

// CreateNoticeResp 创建公告响应结构体
type CreateNoticeResp struct {
	ID       int64 `json:"id"`
	CreateAt int64 `json:"create_at"`
}

// UpdateNoticeReq 更新公告请求结构体
type UpdateNoticeReq struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
}

// UpdateNoticeResp 更新通知响应结构体
type UpdateNoticeResp struct {
	UpdateAt int64 `json:"update_at"`
}

// DeleteNoticeReq 删除公告请求结构体
type DeleteNoticeReq struct {
	ID int64 `json:"id"` // 公告ID
}

// DeleteNoticeResp 删除通知响应结构体
type DeleteNoticeResp struct {
}

// GetNoticeReq 获取通知请求结构体
type GetNoticeReq struct {
	ID int64 `json:"id"` // 公告ID
}

// GetNoticeResp 获取通知响应结构体
type GetNoticeResp struct {
	UpdateAt int64  `json:"update_at"`
	Content  string `json:"content"`
}
