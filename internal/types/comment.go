package types

import "time"

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
