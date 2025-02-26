package controller

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/gin-gonic/gin"
	"strconv"
)

// CollectBlogHandler 收藏帖子
func CollectBlogHandler(c *gin.Context) {
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

	// 尝试收藏帖子
	if err := logic.CollectBlog(ctx, userID, blogID); err != nil {
		if err.Error() == "already collected this blog" {
			zlog.CtxWarnf(ctx, "User already collected this blog (blogID: %d)", blogID)
			response.NewResponse(c).Error(response.PARAM_NOT_VALID) // 已经收藏过该帖子
			return
		}
		zlog.CtxErrorf(ctx, "Failed to collect blog (blogID: %d): %v", blogID, err)
		response.NewResponse(c).Error(response.INTERNAL_ERROR) // 收藏失败
		return
	}

	// 收藏成功
	zlog.CtxInfof(ctx, "Blog collected successfully (blogID: %d)", blogID)
	response.NewResponse(c).Success(nil) // 收藏成功
}

// UncollectBlogHandler 取消收藏帖子
func UncollectBlogHandler(c *gin.Context) {
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

	// 尝试取消收藏帖子
	if err := logic.UncollectBlog(ctx, userID, blogID); err != nil {
		zlog.CtxErrorf(ctx, "Failed to uncollect blog (blogID: %d): %v", blogID, err)
		response.NewResponse(c).Error(response.INTERNAL_ERROR) // 取消收藏失败
		return
	}

	// 取消收藏成功
	zlog.CtxInfof(ctx, "Blog uncollected successfully (blogID: %d)", blogID)
	response.NewResponse(c).Success(nil) // 取消收藏成功
}

// GetCollectedBlogsHandler 获取用户收藏的帖子
func GetCollectedBlogsHandler(c *gin.Context) {
	// 从 Gin 中获取上下文
	ctx := zlog.GetCtxFromGin(c)

	// 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get current user ID: %v", err)
		response.NewResponse(c).Error(response.USER_NOT_LOGIN) // 用户未登录
		return
	}

	// 获取用户收藏的帖子列表
	blogs, err := logic.GetCollectedBlogs(ctx, userID)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get collected blogs for user (userID: %d): %v", userID, err)
		response.NewResponse(c).Error(response.INTERNAL_ERROR) // 获取失败
		return
	}

	// 获取成功
	zlog.CtxInfof(ctx, "Collected blogs retrieved successfully (userID: %d)", userID)
	response.NewResponse(c).Success(gin.H{
		"list": blogs,
	})
}
