package api

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/gin-gonic/gin"
)

// @Title        ai.go
// @Description
// @Create       XdpCs 2025-03-09 下午2:27
// @Update       XdpCs 2025-03-09 下午2:27
func AIChat(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.AIChatReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "Test request: %v", req)
	resp, err := logic.NewAILogic().AIChat(ctx, req)
	response.Response(c, resp, err)
	return
}
