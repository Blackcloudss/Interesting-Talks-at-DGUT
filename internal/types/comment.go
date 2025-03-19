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
