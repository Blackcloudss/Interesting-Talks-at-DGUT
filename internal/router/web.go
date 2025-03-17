package router

import (
	"fmt"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/configs"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/api"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/manager"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/middleware"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
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
	//设置静态路由，用于访问上传的文件
	r.Static("image", "./images")
	// 创建 RouteManager 实例
	routeManager := manager.NewRouteManager(r)
	// 注册各业务路由组的具体路由
	registerRoutes(routeManager)
	return r, nil
}

func registerRoutes(routeManager *manager.RouteManager) {
	//通用功能相关路由
	routeManager.RegisterCommonRoutes(func(rg *gin.RouterGroup) {
		rg.POST("/rtoken", api.RefreshToken) //用rtoken刷新atoken和rtoken
	})

	//微信登陆相关路由
	routeManager.RegisterLoginRoutes(func(rg *gin.RouterGroup) {
		rg.GET("/qrcode", api.GetQRCode)  //获取小程序二维码
		rg.GET("/login", api.WechatLogin) //微信登陆
	})

	//用户信息相关路由
	routeManager.RegisterProfileRoutes(func(rg *gin.RouterGroup) {
		middleware.CheckAtoken()                           // 检查 Atoken
		rg.GET("/phone", api.GetPhone)                     //获取用户手机号（授权时使用）
		rg.GET("/userinfo", api.GetUserInfo)               //获取用户的微信头像和微信昵称（授权时使用）
		middleware.PermissionMiddleware()                  // 检查权限
		rg.GET("/common/show", api.GetCommonProfile)       // 获取用户基本信息
		rg.GET("/private/show", api.GetDetailProfile)      // 获取用户隐私信息
		rg.PUT("/common/update", api.UpdateCommonProfile)  // 更新基本信息
		rg.PUT("/private/update", api.UpdateDetailProfile) // 更新隐私信息
	})
	//帖子相关路由
	routeManager.RegisterBlogRoutes(func(rg *gin.RouterGroup) {
		rg.POST("/create", api.CreateBlogHandler)          // 创建帖子
		rg.PUT("/update", api.UpdateBlogHandler)           // 更新帖子
		rg.DELETE("/delete", api.DeleteBlogHandler)        // 删除帖子
		rg.GET("/get", api.GetBlogByIDHandler)             // 获取帖子详情
		rg.GET("/list", api.GetBlogsHandler)               // 分页显示帖子
		rg.GET("/list_by_tag", api.GetBlogsByTagHandler)   // 根据标签显示帖子列表
		rg.GET("/my_blogs", api.GetMyBlogsHandler)         // 获取当前用户发布的帖子
		rg.POST("/collect", api.CollectBlogHandler)        // 收藏帖子
		rg.POST("/uncollect", api.UncollectBlogHandler)    // 取消收藏帖子
		rg.GET("/collected", api.GetCollectedBlogsHandler) // 获取用户收藏的帖子
		rg.POST("/like", api.LikeBlogHandler)              // 点赞帖子
		rg.POST("/unlike", api.UnlikeBlogHandler)          // 取消点赞
	})

	//评论相关路由
	routeManager.RegisterCommentRoutes(func(rg *gin.RouterGroup) {
		rg.POST("/create", api.CreateComment)   // 创建评论
		rg.DELETE("/delete", api.DeleteComment) // 删除评论
		rg.GET("/list", api.GetCommentList)     // 获取评论列表
	})

}
