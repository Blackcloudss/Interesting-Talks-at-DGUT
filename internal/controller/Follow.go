package controller

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/gin-gonic/gin"
	"strconv"
)

// FollowHandler 关注用户
func FollowHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)

	followerID, err := getCurrentUserID(c)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get current user ID: %v", err)
		response.NewResponse(c).Error(response.USER_NOT_LOGIN) // 用户未登录
		return
	}

	followedID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "Invalid user ID: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID) // 用户ID无效
		return
	}

	if err := logic.Follow(ctx, followerID, followedID); err != nil {
		zlog.CtxErrorf(ctx, "Failed to follow user (followedID: %d): %v", followedID, err)
		response.NewResponse(c).Error(response.INTERNAL_ERROR) // 关注失败
		return
	}

	zlog.CtxInfof(ctx, "Followed successfully (followedID: %d)", followedID)
	response.NewResponse(c).Success(nil) // 关注成功
}

// UnfollowHandler 取消关注用户
func UnfollowHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)

	followerID, err := getCurrentUserID(c)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get current user ID: %v", err)
		response.NewResponse(c).Error(response.USER_NOT_LOGIN) // 用户未登录
		return
	}

	followedID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "Invalid user ID: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID) // 用户ID无效
		return
	}

	if err := logic.Unfollow(ctx, followerID, followedID); err != nil {
		zlog.CtxErrorf(ctx, "Failed to unfollow user (followedID: %d): %v", followedID, err)
		response.NewResponse(c).Error(response.INTERNAL_ERROR) // 取消关注失败
		return
	}

	zlog.CtxInfof(ctx, "Unfollowed successfully (followedID: %d)", followedID)
	response.NewResponse(c).Success(nil) // 取消关注成功
}

// GetFollowingsHandler 获取用户关注的用户列表
func GetFollowingsHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)

	userID, err := getCurrentUserID(c)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get current user ID: %v", err)
		response.NewResponse(c).Error(response.USER_NOT_LOGIN) // 用户未登录
		return
	}

	followings, err := logic.GetFollowings(ctx, userID)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get followings for user (userID: %d): %v", userID, err)
		response.NewResponse(c).Error(response.INTERNAL_ERROR) // 获取失败
		return
	}

	zlog.CtxInfof(ctx, "Followings retrieved successfully (userID: %d)", userID)
	response.NewResponse(c).Success(gin.H{
		"list": followings,
	})
}

// GetFollowersHandler 获取用户的粉丝列表
func GetFollowersHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)

	userID, err := getCurrentUserID(c)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get current user ID: %v", err)
		response.NewResponse(c).Error(response.USER_NOT_LOGIN) // 用户未登录
		return
	}

	followers, err := logic.GetFollowers(ctx, userID)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get followers for user (userID: %d): %v", userID, err)
		response.NewResponse(c).Error(response.INTERNAL_ERROR) // 获取失败
		return
	}

	zlog.CtxInfof(ctx, "Followers retrieved successfully (userID: %d)", userID)
	response.NewResponse(c).Success(gin.H{
		"list": followers,
	})
}
