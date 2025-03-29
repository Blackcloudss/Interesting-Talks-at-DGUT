package types

import (
	"time"
)

// CreateCommentReq 创建评论/回复请求结构体
type CreateCommentReq struct {
	BlogID   int64  `json:"blog_id"`   // 所属帖子 ID
	Content  string `json:"content"`   // 评论内容
	ParentID int64  `json:"parent_id"` // 父评论 ID（等于0为一级回复，不等于0为多级回复）
}

// CreateCommentResp 创建评论响应结构体
type CreateCommentResp struct {
	CommentID int64     `json:"comment_id"` // 创建的评论 ID
	CreatedAt time.Time `json:"created_at"`
}

// DeleteCommentReq 删除评论请求结构体
type DeleteCommentReq struct {
	CommentID int64 `json:"comment_id" binding:"required"`
}

// DeleteCommentResp 删除评论响应结构体
type DeleteCommentResp struct {
}

// LikeCommentReq 点赞/取消点赞评论请求
type LikeCommentReq struct {
	CommentID int64 `json:"comment_id" binding:"required"`
}

// LikeCommentResp 点赞/取消点赞评论响应
type LikeCommentResp struct {
	IsLiked   bool  `json:"is_liked"`   // 操作后的点赞状态
	LikeCount int64 `json:"like_count"` // 操作后的点赞数
}

// CommentDetail 评论详情
type CommentDetail struct {
	ID           int64     `json:"id"`
	BlogID       int64     `json:"blog_id"`
	UserID       int64     `json:"user_id"`
	Content      string    `json:"content"`
	LikeCount    int       `json:"like_count"`
	CreatedAt    time.Time `json:"created_at"`
	Nickname     string    `json:"nickname"` // 回复者昵称
	Avatar       string    `json:"avatar"`   // 回复者头像
	Tag          string    `json:"tag"`      // 回复者标签
	RepliesCount int       `json:"replies_count"`
}

// GetCommentListReq 获取评论列表请求结构体
type GetCommentListReq struct {
	BlogID   int64  ` json:"blog_id" form:"blog_id" `                                     // 所属帖子 ID
	Page     int    `json:"page" form:"page" `                                            // 当前页码
	PageSize int    `json:"page_size" form:"page_size" `                                  // 每页大小
	SortBy   string `json:"sort_by" form:"sort_by" binding:"oneof=created_at like_count"` // 排序方式：created_at（按发布时间）或按like_count（按点赞数）
}

// GetCommentListResp 获取评论列表响应结构体
type GetCommentListResp struct {
	TotalCount int64           `json:"total_count"` // 总评论数
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	Comments   []CommentDetail `json:"comments"` // 评论列表
}

// ReplyDetail 回复详情
type ReplyDetail struct {
	ID           int64     `json:"id"`
	BlogID       int64     `json:"blog_id"`
	UserID       int64     `json:"user_id"`
	Content      string    `json:"content"`
	LikeCount    int       `json:"like_count"`
	CreatedAt    time.Time `json:"created_at"`
	Nickname     string    `json:"nickname"`       // 回复者昵称
	Avatar       string    `json:"avatar"`         // 回复者头像
	Tag          string    `json:"tag"`            // 回复者标签
	ParentUserID int64     `json:"parent_user_id"` // 被回复用户ID
	ParentNick   string    `json:"parent_nick"`    // 被回复用户昵称
}

// GetRepliesListReq 获取回复列表请求
type GetRepliesListReq struct {
	ParentID int64 `form:"parent_id" binding:"required"` // 父评论ID
	Page     int   `form:"page" binding:"required"`      // 页码
	PageSize int   `form:"page_size" binding:"required"` // 每页数量
}

// GetRepliesListResp 回复列表响应
type GetRepliesListResp struct {
	TotalCount int64         `json:"total_count"` // 总回复数
	Page       int           `json:"page"`        // 当前页码
	PageSize   int           `json:"page_size"`   // 每页数量
	Replies    []ReplyDetail `json:"replies"`     // 回复列表
}
