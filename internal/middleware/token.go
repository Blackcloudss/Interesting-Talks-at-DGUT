package middleware

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/jwt"
	"github.com/gin-gonic/gin"
)

// @Title        token.go
// @Description
// @Create       XdpCs 2025-03-10 上午1:18
// @Update       XdpCs 2025-03-10 上午1:18
func ReflashAtoken() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := zlog.GetCtxFromGin(c)
		token := c.GetHeader("Authorization")
		if token == "" {
			zlog.CtxErrorf(ctx, `token is empty`)
			response.NewResponse(c).Error(response.TOKEN_IS_BLANK)
			c.Abort()
			return
		}
		//解析token是否有效，并取出上一次的值
		data, err := jwt.IdentifyToken(ctx, token)
		if err != nil {
			zlog.CtxErrorf(ctx, "ReflashAtoken err:%v", err)
			response.NewResponse(c).Error(response.TOKEN_IS_EXPIRED)
			//对应token无效，直接让他返回
			c.Abort()
			return
		}
		//判断其是否为atoken
		if data.Class != global.AUTH_ENUMS_ATOKEN {
			response.NewResponse(c).Error(response.TOKEN_TYPE_ERROR)
			c.Abort()
			return
		}
		//将token内部数据传下去
		c.Set(global.TOKEN_USER_ID, data.Userid)
		//生成新的token
		resp, err := logic.NewTokenLogic().GenAtoken(ctx, data)
		if err != nil {
			zlog.CtxErrorf(ctx, "ReflashAtoken err:%v", err)
			c.Abort()
			return
		}
		//将token放到响应头
		c.Header("Authorization", resp.Atoken)
		c.Next()
	}
}
