package router

import (
	"Interesting-Talks-at-DGUT/configs"
	"Interesting-Talks-at-DGUT/manager"
	"fmt"
	"github.com/gin-gonic/gin"
)

// @Title        web.go
// @Description
// @Create       XdpCs 2025-02-16 上午1:01
// @Update       XdpCs 2025-02-16 上午1:01

// RunServer 启动服务器 路由层
func RunServer() {
	r, err := listen()
	if err != nil {
		zlog.Errorf("Listen error: %v", err)
		panic(err.Error())
	}
	r.Run(fmt.Sprintf("%s:%d", configs.Conf.App.Host, configs.Conf.App.Port)) // 启动 Gin 服务器
}

// listen 配置 Gin 服务器
func listen() (*gin.Engine, error) {
	r := gin.Default() // 创建默认的 Gin 引擎
	// 注册全局中间件（例如获取 Trace ID）
	manager.RequestGlobalMiddleware(r)
	// 创建 RouteManager 实例
	routeManager := manager.NewRouteManager(r)
	// 注册各业务路由组的具体路由
	registerRoutes(routeManager)
	return r, nil
}

func registerRoutes(routeManager *manager.RouteManager) {
	//通用功能相关路由
	routeManager.RegisterCommonRoutes(func(rg *gin.RouterGroup) {
		rg.POST("/rtoken", api.ReflashRtoken)
	})
}
