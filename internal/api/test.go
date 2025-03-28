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
		return
	}
	zlog.CtxInfof(ctx, "Test request: %v", req)
	resp, err := logic.NewTestLogic().TestLogic(ctx, req)
	response.Response(c, resp, err)
	return
}

// @Title        test.go
func Display(c *gin.Context) {
	const TEST = "" +
		"Hello(╹ڡ╹ ),这里是莞工趣坛项目，看到则说明能访问到服务器。" +
		"本项目还有部分功能尚未完成，敬请期待！" +
		"若后续部分功能无法访问：" +
		"1.请检查下自身的功能url是否正确；" +
		"2.检查参数是否符合文档。" +
		"如有其他问题，请联系qsk❄❄"
	response.Response(c, TEST, nil)
	return
}

// 创建测试成员
func CreateMember(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	resp, err := logic.NewTestLogic().CreateMember(ctx)
	response.Response(c, resp, err)
}
