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
	"strings"
)

//权限校验中间件：检查用户是否有权限访问某个资源 -- 涉及到前端判断是否需要渲染的组件

func PermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := zlog.GetCtxFromGin(c)

		UserId := jwt.GetUserId(c)
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

		// 重要!! 将读取的请求体内容重新设置到 c.Request.Body，供后续处理使用
		c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		// 获取当前请求的URL
		Url = url.GetBaseURL(c)

		// 判断请求方法是否为 DELETE，如果是则截取URL
		if c.Request.Method == "DELETE" {
			u := c.Request.URL.Path
			index := strings.Index(u, "/delete")
			if index != -1 {
				Url = u[:index+len("/delete")]
				zlog.CtxInfof(ctx, "Base URL:%v", Url)
			}
		} else {
			Url = c.Request.URL.Path
			zlog.CtxInfof(ctx, "Base URL:%v", Url)
		}

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
