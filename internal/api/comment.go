package api

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/gin-gonic/gin"
)

// CommentHandler 创建评论
func CreateCommentHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.CreateCommentReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "CreateComment request error: %v", err)
		return
	}
	/*

		// 获取当前用户ID
		userID, err := controller.getCurrentUserID(c)
		if err != nil {
			zlog.CtxErrorf(ctx, "Failed to get current user ID: %v", err)
			response.NewResponse(c).Error(response.USER_NOT_LOGIN) // 用户未登录
			return
		}
		comment.AuthorID = userID*/
	zlog.CtxInfof(ctx, "CreateComment request: %v", req)
	resp, err := logic.NewBlogLogic().CreateComment(ctx, req)
	response.Response(c, resp, err)
	return
}

// DeleteCommentHandler 删除评论
func DeleteCommentHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.DeleteCommentReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "DeleteComment request error: %v", err)
		return
	}
	/*
		userID, err := controller.getCurrentUserID(c)
		if err != nil {
			zlog.CtxErrorf(ctx, "getCurrentUserID failed: %v", err)
			response.NewResponse(c).Error(response.USER_NOT_LOGIN)
			return
		}*/

	zlog.CtxInfof(ctx, "DeleteComment request: %v", req)
	resp, err := logic.NewBlogLogic().DeleteComment(ctx, req)
	response.Response(c, resp, err)
	return
}

// CommentListHandler 获取评论列表
func GetCommentListHandler(c *gin.Context) {
	// 从 Gin 中获取上下文
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetCommentListReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "DeleteComment request error: %v", err)
		return
	}
	/*
		userID, err := controller.getCurrentUserID(c)
		if err != nil {
			zlog.CtxErrorf(ctx, "getCurrentUserID failed: %v", err)
			response.NewResponse(c).Error(response.USER_NOT_LOGIN)
			return
		}*/
	zlog.CtxInfof(ctx, "GetCommentList request: %v", req)
	resp, err := logic.NewBlogLogic().GetCommentList(ctx, req)
	response.Response(c, resp, err)
	return
}
