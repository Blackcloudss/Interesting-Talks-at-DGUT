package middleware

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/repo"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/jwt"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/url"
	"github.com/gin-gonic/gin"
)

const (
	CREATE_POINTS = 2 // 创建积分
)

var INCREASE_POINTS = map[string]int{
	"/api/blog/create":    CREATE_POINTS, // 创建帖子
	"/api/comment/create": CREATE_POINTS, // 创建评论
}

// @Title        point.go
// @Description
// @Create       XdpCs 2025-03-20 上午10:48
// @Update       XdpCs 2025-03-20 上午10:48

// 增加积分
func IncreasePoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := zlog.GetCtxFromGin(c)
		UserId := jwt.GetUserId(c)
		// 获取当前 URL
		Url := url.GetBaseURL(c)
		// 获取当前 URL 对应的积分值
		Point := GetPoints(Url)
		//积分值为0时跳过数据库操作：
		if Point == 0 {
			c.Next()
			return
		}
		err := repo.NewPointRepo(global.DB).IncreasePoints(UserId, Point)
		if err != nil {
			zlog.CtxErrorf(ctx, "IncreasePoints error:%v", err)
			c.Abort()
			return
		}
		c.Next()
	}
}

// 判断当前 URL 符合哪项加分规则
func GetPoints(url string) (points int) {
	points = INCREASE_POINTS[url]
	return
}
