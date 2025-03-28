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
	CommentID int64 `json:"comment_id"` // 要删除的评论 ID
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
type CommentDetail struct {
	ID           int64     `json:"id"`
	BlogID       int64     `json:"blog_id"`
	UserID       int64     `json:"user_id"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	LikesCount   int       `json:"likes_count"`
	RepliesCount int       `json:"replies_count"`
	ParentID     *int64    `json:"parent_id,omitempty"` // 为nil表示一级评论

	// 用户信息
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Tag      string `json:"tag"`

	// 回复列表（可选）
	Replies []CommentDetail `json:"replies,omitempty"`
}

// GetCommentListReq 获取评论列表请求结构体
type GetCommentListReq struct {
	BlogID   int64  `form:"blog_id" binding:"required"`                    // 所属帖子 ID
	Page     int    `form:"page" binding:"required"`                       // 当前页码
	PageSize int    `form:"page_size" binding:"required"`                  // 每页大小
	SortBy   string `form:"sort_by" binding:"oneof=created_at like_count"` // 排序方式：created_at（按发布时间）或按likes_count（按点赞数）
}

// GetCommentListResp 获取评论列表响应结构体
type GetCommentListResp struct {
	TotalCount int64           `json:"total_count"` // 总评论数
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	Comments   []CommentDetail `json:"comments"` // 评论列表
}

// GetRepliesListReq 获取回复列表
type GetRepliesListReq struct {
	ParentID int64 `form:"parent_id" binding:"required"`
	Page     int   `form:"page" binding:"required,min=1"`
	PageSize int   `form:"page_size" binding:"required,min=1,max=100"`
}

/*
// GetSecondCommentListReq 获取更多二级评论请求结构体
type GetSecondCommentListReq struct {
	ParentID int64 `form:"parent_id"`                    // 根评论ID
	Page     int   `form:"page" binding:"required"`      // 当前页码
	PageSize int   `form:"page_size" binding:"required"` // 每页大小
}*/

/*
// GetSecondCommentListResp 获取更多二级评论响应结构体

	type GetSecondCommentListResp struct {
		TotalCount     int64                 `json:"total_count"` // 总评论数
		Page           int                   `json:"page"`
		PageSize       int                   `json:"page_size"`
		SecondComments []SecondCommentDetail `json:"second_comments"` // 二级评论列表
	}
*/

/*
// CommentDetail 评论详情
type CommentDetail struct {
	ID           int64     `json:"id" gorm:"column:id"`
	BlogID       int64     `json:"blog_id" gorm:"column:blog_id"`
	UserID       int64     `json:"user_id" gorm:"column:user_id"`
	Content      string    `json:"content" gorm:"column:content"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
	LikesCount   int       `json:"likes_count" gorm:"column:likes_count"`
	RepliesCount int       `json:"replies_count" gorm:"column:replies_count"`

	// 用户信息
	Nickname string `json:"nickname" gorm:"column:nickname"`
	Avatar   string `json:"avatar" gorm:"column:avatar"`
	Tag      string `json:"tag" gorm:"column:tag"`

	// 二级评论列表（可选）
	SecondComments []SecondCommentDetail `json:"second_comments,omitempty" gorm:"-"`
}

// SecondCommentDetail 二级评论详情
type SecondCommentDetail struct {
	ID        int64     `json:"id" gorm:"column:id"`
	BlogID    int64     `json:"blog_id" gorm:"column:blog_id"`
	UserID    int64     `json:"user_id" gorm:"column:user_id"`
	Content   string    `json:"content" gorm:"column:content"`
	ParentID  int64     `json:"parent_id" gorm:"column:parent_id"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`

	// 用户信息
	Nickname string `json:"nickname" gorm:"column:nickname"`
	Avatar   string `json:"avatar" gorm:"column:avatar"`
	Tag      string `json:"tag" gorm:"column:tag"`
}
*/
