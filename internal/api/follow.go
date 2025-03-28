package api

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/jwt"
	"github.com/gin-gonic/gin"
)

func FollowHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.FollowReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "Follow request error: %v", err)
		return
	}
	zlog.CtxInfof(ctx, "Follow request: %+v", req)
	UserID := jwt.GetUserId(c)
	resp, err := logic.NewFollowLogic().Follow(ctx, req, UserID)
	response.Response(c, resp, err)
}

// GetFollowingsHandler 获取用户关注的用户列表
func GetFollowingsHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetFollowingsReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetFollowings request error: %v", err)
		return
	}
	zlog.CtxInfof(ctx, "GetFollowings request: %+v", req)
	UserID := jwt.GetUserId(c)
	resp, err := logic.NewFollowLogic().GetFollowings(ctx, UserID)
	response.Response(c, resp, err)
}

// GetFollowersHandler 获取用户的粉丝列表
func GetFollowersHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetFollowersReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetFollowers request error: %v", err)
		return
	}
	zlog.CtxInfof(ctx, "GetFollowers request: %+v", req)
	UserID := jwt.GetUserId(c)
	resp, err := logic.NewFollowLogic().GetFollowers(ctx, UserID)
	response.Response(c, resp, err)
}

// GetOtherFollowingsHandler 获取其他用户关注的用户列表
func GetOtherFollowingsHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetOtherFollowingsReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetOtherFollowings request error: %v", err)
		return
	}
	zlog.CtxInfof(ctx, "GetOtherFollowings request: %+v", req)
	resp, err := logic.NewFollowLogic().GetFollowings(ctx, req.OtherUserID)
	response.Response(c, resp, err)
}

// GetOtherFollowersHandler 获取其他用户的粉丝列表
func GetOtherFollowersHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetOtherFollowersReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetOtherFollowers request error: %v", err)
		return
	}
	zlog.CtxInfof(ctx, "GetOtherFollowers request: %+v", req)
	resp, err := logic.NewFollowLogic().GetFollowers(ctx, req.OtherUserID)
	response.Response(c, resp, err)
}
