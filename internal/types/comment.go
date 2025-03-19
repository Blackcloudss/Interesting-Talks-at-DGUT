package types

import "github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"

// CreateCommentReq 创建评论请求结构体
type CreateCommentReq struct {
	BlogID  int64  `json:"blogID"`
	Content string `json:"content"`
}

// CreateCommentResp 创建评论响应结构体
type CreateCommentResp struct {
	CommentID int64 `json:"commentID"`
	CreatedAt int64 `json:"createdAt"`
}

// DeleteCommentReq 删除评论请求结构体
type DeleteCommentReq struct {
	CommentID int64 `json:"commentID"`
}

// DeleteCommentResp 删除评论响应结构体
type DeleteCommentResp struct {
}

// GetCommentListReq 获取评论列表请求体
type GetCommentListReq struct {
	BlogID int64 `json:"blogID"`
}

// GetCommentListResp 获取评论列表响应体
type GetCommentListResp struct {
	Comments []model.Comment `json:"comments"`
}

// LikeCommentReq 点赞评论请求结构体
type LikeCommentReq struct {
	CommentID int64 `json:"commentID"`
}

// LikeCommentResp 点赞评论响应结构体
type LikeCommentResp struct {
}

// UnlikeCommentReq 取消点赞评论请求结构体
type UnlikeCommentReq struct {
	CommentID int64 `json:"commentID"`
}

// UnlikeCommentResp 取消点赞评论响应结构体
type UnlikeCommentResp struct {
}

// ReplyCommentReq 回复评论请求结构体
type ReplyCommentReq struct {
	ParentID int64  `json:"parentID"` // 父评论 ID
	Content  string `json:"content"`  // 回复内容
}

// ReplyCommentResp 回复评论响应结构体
type ReplyCommentResp struct {
	CommentID int64 `json:"commentID"` // 回复评论的 ID
	CreatedAt int64 `json:"createdAt"` // 创建时间戳
}

// GetRepliesListReq 获取回复列表请求结构体
type GetRepliesListReq struct {
	CommentID int64 `json:"commentID"` // 父评论 ID
}

// GetRepliesListResp 获取回复列表响应结构体
type GetRepliesListResp struct {
	Replies []model.Comment `json:"replies"` // 回复列表
}
