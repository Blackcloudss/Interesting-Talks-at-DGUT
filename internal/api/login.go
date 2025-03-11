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
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "WechatLogin request: %v", req)
	resp, err := logic.NewWechatLoginLogic().WechatLogin(ctx, req)
	response.Response(c, resp, err)
	return
}
