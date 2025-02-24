package manager

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/middleware"
	"github.com/gin-gonic/gin"
)

// @Title        manager.go
// @Description
// @Create       XdpCs 2025-02-24 下午4:03
// @Update       XdpCs 2025-02-24 下午4:03

// PathHandler 是一个用于注册路由组的函数类型
type PathHandler func(rg *gin.RouterGroup)

// Middleware 是一个用于生成中间件的函数类型
type Middleware func() gin.HandlerFunc

// RouteManager 管理不同的路由组，按业务功能分组
type RouteManager struct {
}

// NewRouteManager 创建一个新的 RouteManager 实例，包含各业务功能的路由组
func NewRouteManager(router *gin.Engine) *RouteManager {
	return &RouteManager{}
}

// RegisterMiddleware 根据组名为对应的路由组注册中间件
// group 参数为 "login"、"profile"、"team"或"Common"，分别对应不同的路由组
func (rm *RouteManager) RegisterMiddleware(group string, middleware Middleware) {
	switch group {
	}
}

// RequestGlobalMiddleware 注册全局中间件，应用于所有路由
func RequestGlobalMiddleware(r *gin.Engine) {
	r.Use(requestid.New())
	r.Use(middleware.AddTraceId())
	r.Use(middleware.Cors())
}
