package initalize

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"runtime"
)

// @Title        eve.go
// @Description
// @Create       XdpCs 2025-02-16 上午12:55
// @Update       XdpCs 2025-02-16 上午12:55
func Eve() {
	zlog.Warnf("开始释放资源！")
	errRedis := global.Rdb.Close()
	if errRedis != nil {
		zlog.Errorf("Redis关闭失败 ：%v", errRedis.Error())
	}
	sqlDB, _ := global.DB.DB()
	errDB := sqlDB.Close()
	if errDB != nil {
		zlog.Errorf("数据库关闭失败 ：%v", errDB.Error())
	}
	runtime.GC()
	if errDB == nil && errRedis == nil {
		zlog.Warnf("资源释放成功！")
	}
}
