package api

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/gin-gonic/gin"
)

// @Title        token.go
// @Description
// @Create       XdpCs 2025-03-10 下午4:00
// @Update       XdpCs 2025-03-10 下午4:00
func RefreshToken(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.TokenReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "RefreshRtoken request error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "RefreshRtoken request: %v", req)
	resp, err := logic.NewTokenLogic().RefreshToken(ctx, req)
	response.Response(c, resp, err)
}
