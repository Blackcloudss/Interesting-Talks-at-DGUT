package types

import "github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"

// CreateCommentReq 创建评论请求结构体
type CreateCommentReq struct {
	BlogID  int64  `json:"blogID"`
	UserID  int64  `json:"userID"`
	Content string `json:"content"`
}

// CreateCommentResp 创建评论响应结构体
type CreateCommentResp struct {
	Comment model.Comment `json:"comment"`
}

// DeleteCommentReq 删除评论请求结构体
type DeleteCommentReq struct {
	CommentID int64 `json:"commentID"`
	UserID    int64 `json:"userID"`
}

// DeleteCommentResp 删除评论响应结构体
type DeleteCommentResp struct {
	Success bool `json:"success"`
}

// 获取评论列表请求体
type GetCommentListReq struct {
	BlogID int64 `json:"blogID"`
}

// 获取评论列表响应体
type GetCommentResp struct {
	Comment []model.Comment `json:"comment"`
}
