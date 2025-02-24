package initalize

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/configs"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
)

// @Title        init_log.go
// @Description
// @Create       XdpCs 2025-02-16 上午12:58
// @Update       XdpCs 2025-02-16 上午12:58
func InitLog(config *configs.Config) {
	logger := log.GetZap(config)
	zlog.InitLogger(logger)
}
