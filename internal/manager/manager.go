package manager

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/middleware"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

// @Title        manager.go
// @Description
// @Create       XdpCs 2025-02-24 下午4:03
// @Update       XdpCs 2025-02-24 下午4:03

// PathHandler 是一个用于注册路由组的函数类型
type PathHandler func(rg *gin.RouterGroup)

// RouteManager 管理不同的路由组，按业务功能分组
type RouteManager struct {
	CommonRoutes  *gin.RouterGroup // 通用的路由组
	LoginRoutes   *gin.RouterGroup // 登录相关的路由组
	ProfileRoutes *gin.RouterGroup // 个人信息相关路由组
	AIRoutes      *gin.RouterGroup // AI相关路由组
}

// NewRouteManager 创建一个新的 RouteManager 实例，包含各业务功能的路由组
func NewRouteManager(router *gin.Engine) *RouteManager {
	return &RouteManager{
		CommonRoutes:  router.Group("/api/common"),  // 初始化通用路由组
		LoginRoutes:   router.Group("/api/wxlogin"), // 初始化登录路由组
		ProfileRoutes: router.Group("/api/profile"), // 初始化个人信息路由组
		AIRoutes:      router.Group("/api/ai"),      //初始化AI路由组
	}
}

func (rm *RouteManager) RegisterCommonRoutes(handler PathHandler) {
	handler(rm.LoginRoutes)
}

func (rm *RouteManager) RegisterLoginRoutes(handler PathHandler) {
	handler(rm.LoginRoutes)
}

func (rm *RouteManager) RegisterProfileRoutes(handler PathHandler) {
	handler(rm.ProfileRoutes)
}

// RequestGlobalMiddleware 注册全局中间件，应用于所有路由
func RequestGlobalMiddleware(r *gin.Engine) {
	//为每个请求生成唯一的请求ID，方便追踪和日志记录
	r.Use(requestid.New())
	r.Use(middleware.AddTraceId())
	r.Use(middleware.Cors())
}
