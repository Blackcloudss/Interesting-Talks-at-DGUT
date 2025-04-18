package api

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/gin-gonic/gin"
)

// @Title        login.go
// @Description
// @Create       XdpCs 2025-03-06 上午1:23
// @Update       XdpCs 2025-03-06 上午1:23
func WechatLogin(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.WechatLoginReq](c) // 修改请求结构体
	if err != nil {
		zlog.CtxErrorf(ctx, "WechatLogin request error: %v", err)
		return
	}
	zlog.CtxInfof(ctx, "WechatLogin request: %v", req)
	resp, err := logic.NewWechatLoginLogic().WechatLogin(ctx, req)

	zlog.CtxErrorf(ctx, "微信登陆错误打印: %v", err)

	response.Response(c, resp, err)
	return
}

// 获取小程序二维码
func GetQRCode(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.QRCodeReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetQRCode request error: %v", err)
		return
	}
	zlog.CtxInfof(ctx, "GetQRCode request: %v", req)
	resp, err := logic.NewWechatLoginLogic().GetQRCode(ctx)
	response.Response(c, resp, err)
}
