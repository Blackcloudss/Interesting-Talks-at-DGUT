package url

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/gin-gonic/gin"
	"strings"
)

// @Title        url.go
// @Description
// @Create       XdpCs 2025-03-20 上午11:55
// @Update       XdpCs 2025-03-20 上午11:55

func GetBaseURL(c *gin.Context) (Url string) {
	ctx := zlog.GetCtxFromGin(c)
	// 如果是delete请求，需要截取 url的前面部分
	if c.Request.Method == "DELETE" {
		url := c.Request.URL.Path
		index := strings.Index(url, "/delete")
		if index != -1 {
			// 切片获取，左闭右开
			Url = url[:index+len("/delete")]
			zlog.CtxInfof(ctx, "Base URL:%v", Url)
			return Url
		}
	} else {
		Url = c.Request.URL.Path
		zlog.CtxInfof(ctx, "Base URL:%v", Url)
		return Url
	}
	return
}
