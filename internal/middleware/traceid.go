package middleware

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Title        traceid.go
// @Description
// @Create       XdpCs 2025-02-24 下午3:58
// @Update       XdpCs 2025-02-24 下午3:58
func AddTraceId() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头中的 Trace ID
		traceID := c.GetHeader("X-Request-ID")
		if traceID == "" {
			// 如果没有 Trace ID，生成一个新的,  默认节点为 1
			node, _ := snowflake.NewNode(global.DEFAULT_NODE_ID)
			traceID = node.Generate().String()
		}

		// 将 Trace ID 存入上下文中，方便后续处理使用
		zlog.SetCtxFromGin(c, zlog.NewContext(c.Request.Context(), zap.String(zlog.LogKeyTraceId, traceID)))
		// 在日志中记录 Trace ID
		zlog.CtxInfof(c, "TraceID: %s", traceID)

		// 继续执行下一个中间件或请求处理
		c.Next()
	}
}
