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
	// always in the first
	ctx := zlog.GetCtxFromGin(c)

	req, err := types.BindReq[types.TestO1Req](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "Test request: %v", req)
	resp, err := logic.NewTestLogic().TestLogic(ctx, req)
	// 更加人性化的response返回，这样减少重复代码的书写
	response.Response(c, resp, err)

	return
}
