package api

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/jwt"
	"github.com/gin-gonic/gin"
)

// @Title        phone.go
// @Description
// @Create       XdpCs 2025-03-11 上午9:29
// @Update       XdpCs 2025-03-11 上午9:29
func GetPhone(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	userid := jwt.GetUserId(c)
	// 获取 code
	req, err := types.BindReq[types.WxPhoneReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetPhone error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "GetPhone request: %v", req)
	resp, err := logic.NewPhoneLogic().GetPhone(ctx, req, userid)
	response.Response(c, resp, err)
	return
}
