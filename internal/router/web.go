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
	err = r.Run(fmt.Sprintf("%s:%d", configs.Conf.App.Host, configs.Conf.App.Port))
	if err != nil {
		zlog.Errorf("RunServer error: %v", err)
		return
	} // 启动 Gin 服务器
}

// listen 配置 Gin 服务器
func listen() (*gin.Engine, error) {
	r := gin.Default() // 创建默认的 Gin 引擎
	// 注册全局中间件（例如获取 Trace ID）
	manager.RequestGlobalMiddleware(r)
	//设置静态路由，用于访问上传的文件
	r.Static("images", "./images")
	// 创建 RouteManager 实例
	routeManager := manager.NewRouteManager(r)
	// 注册各业务路由组的具体路由
	registerRoutes(routeManager)
	return r, nil
}

func registerRoutes(routeManager *manager.RouteManager) {
	//测试相关路由
	routeManager.RegisterTestRoutes(func(rg *gin.RouterGroup) {
		rg.GET("", api.Display) // 测试路由
	})

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
		rg.Use(middleware.CheckAtoken())      // 检查 Atoken
		rg.PUT("/userinfo", api.SaveUserInfo) // 保存用户的微信头像和微信昵称（授权时使用）
		rg.GET("/phone", api.GetPhone)        // 获取用户手机号（授权时使用）
		CommonProfile := rg.Group("/common")
		{
			CommonProfile.GET("/show", api.GetCommonProfile)      // 获取用户基本信息
			CommonProfile.GET("/other", api.GetOtherProfile)      // 获取他人基本信息
			CommonProfile.PUT("/update", api.UpdateCommonProfile) // 更新基本信息
		}
		PrivateProfile := rg.Group("/private")
		{
			PrivateProfile.GET("/show", api.GetPrivateProfile)      // 获取用户隐私信息
			PrivateProfile.PUT("/update", api.UpdatePrivateProfile) // 更新隐私信息
		}
		rg.Use(middleware.PermissionMiddleware()) // 验证权限
		rg.PUT("/role", api.UpdateOtherRole)      // 修改其他用户角色
	})

	//聊天相关路由
	routeManager.RegisterChatRoutes(func(rg *gin.RouterGroup) {
		rg.Use(middleware.CheckAtoken())           // 检查 Atoken
		rg.GET("/friends", api.GetFriendList)      // 获取好友列表
		rg.GET("/messages", api.GetHistoryMessage) // 获取聊天记录
		rg.GET("/wschat", api.WebSocketHandler)    // 建立WebSocket连接，和指定好友聊天
	})

	//AI相关路由
	routeManager.RegisterAIRoutes(func(rg *gin.RouterGroup) {
		rg.Use(middleware.CheckAtoken())                      // 检查 Atoken
		rg.POST("/chat", api.AIChatStream)                    // AI流式聊天
		rg.DELETE("/delete/:type", api.DeleteHistoryMessages) // 删除AI历史聊天记录
	})

	//帖子相关路由
	routeManager.RegisterBlogRoutes(func(rg *gin.RouterGroup) {
		rg.Use(middleware.CheckAtoken())                                      // 检查 Atoken
		rg.POST("/create", middleware.IncreasePoint(), api.CreateBlogHandler) // 创建帖子
		rg.PUT("/update", api.UpdateBlogHandler)                              // 更新帖子
		rg.DELETE("/delete/:blog_id", api.DeleteBlogHandler)                  // 删除帖子
		rg.GET("/get", api.GetBlogByIDHandler)                                // 获取帖子详情
		rg.GET("/list", api.GetBlogsHandler)                                  // 首页显示帖子
		rg.GET("/list_by_tag", api.GetBlogsByTagHandler)                      // 根据标签显示帖子列表
		rg.GET("/my_blogs", api.GetMyBlogsHandler)                            // 获取当前用户发布的帖子
		rg.GET("/other_blogs", api.GetBlogsByUserIDHandler)                   //获取其他用户的帖子（点击其他用户主页可看到）
		rg.POST("/collect", api.CollectBlogHandler)                           // 收藏帖子
		rg.POST("/uncollect", api.UncollectBlogHandler)                       // 取消收藏帖子
		rg.GET("/collected", api.GetCollectedBlogsHandler)                    // 获取用户收藏的帖子（我的帖子）
		rg.POST("/like", api.LikeBlogHandler)                                 // 点赞帖子
		rg.POST("/unlike", api.UnlikeBlogHandler)                             // 取消点赞
	})

	//搜索模块相关路由
	routeManager.RegisterSearchRoutes(func(rg *gin.RouterGroup) {
		rg.Use(middleware.CheckAtoken())                           // 检查 Atoken
		rg.GET("/blogs", api.SearchBlogsHandler)                   //关键词搜索帖子
		rg.GET("/get_history", api.GetSearchHistoryHandler)        //获取搜索历史
		rg.POST("/delete_history", api.DeleteSearchHistoryHandler) //删除搜索历史
		rg.GET("/hot_search", api.GetHotSearch)                    //热搜榜
	})

	//评论相关路由
	routeManager.RegisterCommentRoutes(func(rg *gin.RouterGroup) {
		rg.Use(middleware.CheckAtoken())                                  // 检查 Atoken
		rg.POST("/create", middleware.IncreasePoint(), api.CreateComment) // 创建评论
		rg.DELETE("/delete/:comment_id", api.DeleteComment)               // 删除评论
		rg.GET("/list", api.GetCommentList)                               // 获取一级评论列表（显示部分二级评论）
		rg.GET("/secondlist", api.GetMoreSecondComment)                   //获取更多二级评论
		rg.POST("/like", api.LikeComment)                                 //点赞评论
		rg.POST("/unlike", api.UnlikeComment)                             //取消点赞评论
	})

	//公告相关路由
	routeManager.RegisterNoticeRoutes(func(rg *gin.RouterGroup) {
		rg.Use(middleware.CheckAtoken())                         // 检查 Atoken
		rg.Use(middleware.PermissionMiddleware())                // 验证权限
		rg.POST("/create", api.CreateNoticeHandler)              //创建公告
		rg.PUT("/update", api.UpdateNoticeHandler)               //修改公告
		rg.DELETE("/delete/:notice_id", api.DeleteNoticeHandler) //删除公告
		rg.GET("/get", api.GetNoticeHandler)                     //获取公告
	})

	//关注相关路由
	routeManager.RegisterFollowRoutes(func(rg *gin.RouterGroup) {
		rg.Use(middleware.CheckAtoken())                           // 检查 Atoken
		rg.POST("/follow", api.FollowHandler)                      // 关注用户
		rg.POST("/unfollow", api.UnfollowHandler)                  // 取消关注用户
		rg.GET("/followings", api.GetFollowingsHandler)            // 获取当前用户关注的用户列表
		rg.GET("/followers", api.GetFollowersHandler)              // 获取当前用户的粉丝列表
		rg.GET("/followings/other", api.GetOtherFollowingsHandler) // 获取其他用户关注的用户列表
		rg.GET("/followers/other", api.GetOtherFollowersHandler)   // 获取其他用户的粉丝列表
	})

}
