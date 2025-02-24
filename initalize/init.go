package initalize

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils"
)

// @Title        init.go
// @Description
// @Create       XdpCs 2025-02-01 下午7:59
// @Update       XdpCs 2025-02-01 下午7:59
func Init() {
	introduce()
	InitLog(global.Config)
	InitPath()
	InitConfig()
	InitLog(global.Config)
	InitDataBase(*global.Config)
	InitRedis(*global.Config)
}

func InitPath() {
	global.Path = utils.GetRootPath("")
}
