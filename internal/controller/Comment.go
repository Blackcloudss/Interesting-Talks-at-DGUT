package controller

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/gin-gonic/gin"
	"strconv"
)

// CommentHandler 创建评论
func CommentHandler(c *gin.Context) {
	// 从 Gin 中获取上下文
	ctx := zlog.GetCtxFromGin(c)

	var comment model.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		zlog.CtxErrorf(ctx, "Failed to bind JSON data: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID) // 参数绑定失败，返回参数无效错误
		return
	}

	// 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get current user ID: %v", err)
		response.NewResponse(c).Error(response.USER_NOT_LOGIN) // 用户未登录
		return
	}
	comment.AuthorID = userID

	//创建评论
	if err := logic.CreateComment(ctx, &comment); err != nil {
		zlog.CtxErrorf(ctx, "Failed to create comment: %v", err)
		response.NewResponse(c).Error(response.INTERNAL_ERROR) // 创建评论失败
		return
	}

	zlog.CtxInfof(ctx, "Comment created successfully: %+v", comment)
	response.NewResponse(c).Success(gin.H{ // 返回评论成功响应
		"commentID": comment.CommentID,
	})
}

// DeleteCommentHandler 删除评论
func DeleteCommentHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "ParseUint failed: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}

	userID, err := getCurrentUserID(c)
	if err != nil {
		zlog.CtxErrorf(ctx, "getCurrentUserID failed: %v", err)
		response.NewResponse(c).Error(response.USER_NOT_LOGIN)
		return
	}

	if err := logic.DeleteComment(ctx, commentID, userID); err != nil {
		zlog.CtxErrorf(ctx, "DeleteComment failed: %v", err)
		response.NewResponse(c).Error(response.INTERNAL_ERROR)
		return
	}

	response.NewResponse(c).Success(nil)
}

// CommentListHandler 获取评论列表
func CommentListHandler(c *gin.Context) {
	// 从 Gin 中获取上下文
	ctx := zlog.GetCtxFromGin(c)

	blogIDStr := c.Query("blogID")
	if blogIDStr == "" {
		zlog.CtxErrorf(ctx, "blogID is required")
		response.NewResponse(c).Error(response.PARAM_IS_BLANK) // blogID为空
		return
	}

	blogID, err := strconv.ParseUint(blogIDStr, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "Invalid blogID: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID) // blogID无效
		return
	}

	// 获取评论列表
	comments, err := logic.GetCommentList(ctx, blogID)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to retrieve comment list: %v", err)
		response.NewResponse(c).Error(response.INTERNAL_ERROR) // 获取评论失败
		return
	}

	zlog.CtxInfof(ctx, "Comment list retrieved successfully: %+v", comments)
	response.NewResponse(c).Success(comments) // 返回评论列表
}
