package controller

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/gin-gonic/gin"
	"strconv"
)

// LikeBlogHandler 点赞帖子
func LikeBlogHandler(c *gin.Context) {
	// 从 Gin 中获取上下文
	ctx := zlog.GetCtxFromGin(c)

	// 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get current user ID: %v", err)
		response.NewResponse(c).Error(response.USER_NOT_LOGIN) // 用户未登录
		return
	}

	// 获取帖子ID并校验
	blogID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "Invalid blog ID: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID) // 帖子ID无效
		return
	}

	//点赞帖子
	if err := logic.LikeBlog(ctx, userID, blogID); err != nil {
		if err.Error() == "already liked this blog" {
			zlog.CtxWarnf(ctx, "User already liked this blog (blogID: %d)", blogID)
			response.NewResponse(c).Error(response.PARAM_NOT_VALID) // 已经点赞过该帖子
			return
		}
		zlog.CtxErrorf(ctx, "Failed to like blog (blogID: %d): %v", blogID, err)
		response.NewResponse(c).Error(response.INTERNAL_ERROR) // 点赞失败
		return
	}

	// 点赞成功
	zlog.CtxInfof(ctx, "Blog liked successfully (blogID: %d)", blogID)
	response.NewResponse(c).Success(nil) // 点赞成功
}

// UnlikeBlogHandler 取消点赞
func UnlikeBlogHandler(c *gin.Context) {
	// 从 Gin 中获取上下文
	ctx := zlog.GetCtxFromGin(c)

	// 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get current user ID: %v", err)
		response.NewResponse(c).Error(response.USER_NOT_LOGIN) // 用户未登录
		return
	}

	// 获取帖子ID并校验
	blogID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "Invalid blog ID: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID) // 帖子ID无效
		return
	}

	//取消点赞
	if err := logic.UnlikeBlog(ctx, userID, blogID); err != nil {
		zlog.CtxErrorf(ctx, "Failed to unlike blog (blogID: %d): %v", blogID, err)
		response.NewResponse(c).Error(response.INTERNAL_ERROR) // 取消点赞失败
		return
	}

	// 取消点赞成功
	zlog.CtxInfof(ctx, "Blog unliked successfully (blogID: %d)", blogID)
	response.NewResponse(c).Success(nil) // 取消点赞成功
}
