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
	userid := jwt.GetUserId(c)
	req, err := types.BindReq[types.GetCommonProfileReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetCommonProfile error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "GetCommonProfile request: %v", req)
	resp, err := logic.NewUserLogic().GetCommonProfile(ctx, userid)
	response.Response(c, resp, err)
}

// 获取其他用户信息
func GetOtherProfile(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetOtherProfileReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetOtherProfile error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "GetOtherProfile request: %v", req)
	OtherID := req.OtherID
	resp, err := logic.NewUserLogic().GetCommonProfile(ctx, OtherID)
	response.Response(c, resp, err)
}

// 更新用户基本信息
func UpdateCommonProfile(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	userid := jwt.GetUserId(c)
	req, err := types.BindReq[types.UpdateCommonProfileReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "UpdateCommonProfile error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "UpdateCommonProfile request: %v", req)
	resp, err := logic.NewUserLogic().UpdateCommonProfile(ctx, userid, req)
	response.Response(c, resp, err)
}

// 获取用户隐私信息
func GetPrivateProfile(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	userid := jwt.GetUserId(c)
	req, err := types.BindReq[types.GetPrivateProfileReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetDetailProfile error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "GetDetailProfile request: %v", req)
	resp, err := logic.NewUserLogic().GetPrivateProfile(ctx, userid)
	response.Response(c, resp, err)
}

// 更新用户隐私信息
func UpdatePrivateProfile(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	userid := jwt.GetUserId(c)
	req, err := types.BindReq[types.UpdatePrivateProfileReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "UpdateDetailProfile error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "UpdateDetailProfile request: %v", req)
	resp, err := logic.NewUserLogic().UpdatePrivateProfile(ctx, userid, req)
	response.Response(c, resp, err)
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
