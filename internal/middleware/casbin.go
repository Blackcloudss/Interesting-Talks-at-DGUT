package middleware

import (
	"bytes"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/repo"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/jwt"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/url"
	"github.com/gin-gonic/gin"
	"io"
)

//权限校验中间件：检查用户是否有权限访问某个资源 -- 涉及到前端判断是否需要渲染的组件

func PermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := zlog.GetCtxFromGin(c)

		UserId := jwt.GetUserId(c) //正式使用，
		//var UserId int64          测试使用
		//var req any               测试使用
		var Url string
		var err error

		// 在请求处理中读取并缓存请求体
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			zlog.CtxErrorf(ctx, "读取请求体失败: %v", err)
		} else {
			// 保存请求体内容供后续使用
			//使用 cachedRequestBody 保存的内容来做日志记录或其它操作
			cachedRequestBody := string(bodyBytes)
			zlog.CtxInfof(ctx, "请求体内容: %s", cachedRequestBody)
			// 将读取的请求体内容重新设置到 c.Request.Body，供后续处理使用
			c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		//req, err = types.BindReq[types.RuleCheck](c)
		//if err != nil {
		//	zlog.CtxErrorf(ctx, "PermissionMiddleware 参数绑定失败: %v", err)
		//	zlog.CtxInfof(ctx, "请求详情: Method=%s, Headers=%v, Query=%v", c.Request.Method, c.Request.Header, c.Request.URL.Query())
		//	response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		//	c.Abort()
		//	return
		//}
		//zlog.CtxInfof(ctx, "PermissionMiddleware middleware: %v", req)

		// 重要!! 将读取的请求体内容重新设置到 c.Request.Body，供后续处理使用
		c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		// 获取当前请求的URL
		Url = url.GetBaseURL(c)

		// CheckUserPermissions 检查用户权限
		exist, err := repo.NewCasbinRepo(global.DB).CheckUserPermission(Url, UserId)
		if err != nil {
			response.NewResponse(c).Error(response.PARAM_NOT_VALID)
			c.Abort()
			return
		}
		if exist == false {
			response.NewResponse(c).Error(response.INSUFFICENT_PERMISSIONS)
			c.Abort()
			return
		}
		c.Next() // 继续处理请求
	}
}
