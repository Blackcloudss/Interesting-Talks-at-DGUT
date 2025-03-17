package api

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/jwt"
	"github.com/gin-gonic/gin"
)

// @Title        tz_user.go
// @Description
// @Create       XdpCs 2025-03-10 下午4:40
// @Update       XdpCs 2025-03-10 下午4:40

// 获取用户基本信息
func GetCommonProfile(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.CommonProfileReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetCommonProfile error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "GetCommonProfile request: %v", req)
}

// 更新用户基本信息
func UpdateCommonProfile(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.CommonProfileReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "UpdateCommonProfile error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "UpdateCommonProfile request: %v", req)
}

// 获取用户隐私信息
func GetDetailProfile(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.PrivateProfileReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetDetailProfile error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "GetDetailProfile request: %v", req)
}

// 更新用户隐私信息
func UpdateDetailProfile(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.CommonProfileReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "UpdateDetailProfile error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "UpdateDetailProfile request: %v", req)
}

// 获取用户的微信头像和微信昵称（授权时使用）
func GetUserInfo(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	userid := jwt.GetUserId(c)
	// 获取 用户加密信息
	req, err := types.BindReq[types.UserInfoReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetUserInfo error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "GetUserInfo request: %v", req)
	resp, err := logic.NewUserLogic().GetUserInfo(ctx, userid, req)
	response.Response(c, resp, err)
}
