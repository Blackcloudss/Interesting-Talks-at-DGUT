package api

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/gin-gonic/gin"
)

// FollowHandler 关注用户
func FollowHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.FollowReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "Follow request error: %v", err)
		response.Response(c, nil, err)
		return
	}

	zlog.CtxInfof(ctx, "Follow request: %v", req)
	err = logic.NewFollowLogic().Follow(ctx, req)
	response.Response(c, nil, err)
}

// UnfollowHandler 取消关注用户
func UnfollowHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.UnfollowReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "Unfollow request error: %v", err)
		response.Response(c, nil, err)
		return
	}

	zlog.CtxInfof(ctx, "Unfollow request: %v", req)
	err = logic.NewFollowLogic().Unfollow(ctx, req)
	response.Response(c, nil, err)
}

// GetFollowingsHandler 获取用户关注的用户列表
func GetFollowingsHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetFollowingsReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetFollowings request error: %v", err)
		response.Response(c, nil, err)
		return
	}

	zlog.CtxInfof(ctx, "GetFollowings request: %v", req)
	resp, err := logic.NewFollowLogic().GetFollowings(ctx, req)
	response.Response(c, resp, err)
}

// GetFollowersHandler 获取用户的粉丝列表
func GetFollowersHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetFollowersReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetFollowers request error: %v", err)
		response.Response(c, nil, err)
		return
	}
	zlog.CtxInfof(ctx, "GetFollowers request: %v", req)
	resp, err := logic.NewFollowLogic().GetFollowers(ctx, req)
	response.Response(c, resp, err)
}
