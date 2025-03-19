package api

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/gin-gonic/gin"
)

// @Title        test.go
// @Description
// @Create       XdpCs 2025-02-24 下午11:22
// @Update       XdpCs 2025-02-24 下午11:22
// Test  api层 仅作为校验参数和返回相应，复杂逻辑交给logic层处理
func Test(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)

	req, err := types.BindReq[types.TestO1Req](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "Test request error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "Test request: %v", req)
	resp, err := logic.NewTestLogic().TestLogic(ctx, req)
	response.Response(c, resp, err)
	return
}
