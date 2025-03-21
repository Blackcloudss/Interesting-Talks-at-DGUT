package types

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"time"
)

// CreateCommentReq 创建评论/回复请求结构体
type CreateCommentReq struct {
	BlogID       int64  `json:"blog_id"`   // 所属帖子 ID
	Content      string `json:"content"`   // 评论内容
	ParentID     int64  `json:"parent_id"` // 父评论 ID（等于0为根评论的回复，不等于0为回复的回复）
	RootParentID int64  `json:"origin_id"` // 根评论ID（=0为帖子的直接评论——一级评论）
}

// CreateCommentResp 创建评论响应结构体
type CreateCommentResp struct {
	CommentID int64     `json:"comment_id"` // 创建的评论 ID
	CreatedAt time.Time `json:"created_at"`
}

// DeleteCommentReq 删除评论请求结构体
type DeleteCommentReq struct {
	CommentID int64 `json:"commentID"` // 要删除的评论 ID
}

// DeleteCommentResp 删除评论响应结构体
type DeleteCommentResp struct {
}

// LikeCommentReq 点赞评论请求结构体
type LikeCommentReq struct {
	CommentID int64 `json:"commentID"` // 要点赞的评论 ID
}

// LikeCommentResp 点赞评论响应结构体
type LikeCommentResp struct {
}

// UnlikeCommentReq 取消点赞评论请求结构体
type UnlikeCommentReq struct {
	CommentID int64 `json:"commentID"` // 要取消点赞的评论 ID
}

// UnlikeCommentResp 取消点赞评论响应结构体
type UnlikeCommentResp struct {
}

// GetCommentListReq 获取评论列表请求结构体
type GetCommentListReq struct {
	BlogID   int64  `json:"blog_id" binding:"required"`                     // 所属帖子 ID
	Page     int    `json:"page" binding:"required"`                        // 当前页码
	PageSize int    `json:"page_size" binding:"required"`                   // 每页大小
	SortBy   string `json:"sort_by" binding:"oneof=created_at likes_count"` // 排序方式：created_at（按发布时间）或按likes_count（按点赞数）
}

// GetCommentListResp 获取评论列表响应结构体
type GetCommentListResp struct {
	TotalCount int64           `json:"total_count"` // 总评论数
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	Comments   []CommentDetail `json:"comments"` // 评论列表
}

// GetSecondCommentListReq 获取更多二级评论请求结构体
type GetSecondCommentListReq struct {
	RootParentID int64 `json:"root_parent_id" binding:"required"` // 根评论ID
	Page         int   `json:"page" binding:"required"`           // 当前页码
	PageSize     int   `json:"page_size" binding:"required"`      // 每页大小
}

// GetSecondCommentListResp 获取更多二级评论响应结构体
type GetSecondCommentListResp struct {
	TotalCount     int64                 `json:"total_count"` // 总评论数
	Page           int                   `json:"page"`
	PageSize       int                   `json:"page_size"`
	SecondComments []SecondCommentDetail `json:"second_comments"` // 二级评论列表
}

// CommentDetail 评论详情
type CommentDetail struct {
	FirstComment   model.FirstComment    `json:"first_comment"`
	Nickname       string                `json:"nickname"`        // 评论者昵称
	Avatar         string                `json:"avatar"`          // 评论者头像
	Tag            string                `json:"tag"`             // 评论者标签
	SecondComments []SecondCommentDetail `json:"second_comments"` // 二级评论
}

// SecondCommentDetail 二级评论详情
type SecondCommentDetail struct {
	SecondComment model.SecondComment `json:"second_comment"`
	Nickname      string              `json:"nickname"` // 评论者昵称
	Avatar        string              `json:"avatar"`   // 评论者头像
	Tag           string              `json:"tag"`      // 评论者标签
}
