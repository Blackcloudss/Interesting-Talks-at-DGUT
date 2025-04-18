package main

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/initalize"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/router"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
)

// @Title        main.go
// @Description
// @Create       XdpCs 2025-02-01 下午7:10
// @Update       XdpCs 2025-02-01 下午7:10
func main() {
	initalize.Init()
	defer initalize.Eve()
	router.RunServer()
	zlog.Infof("程序运行完成！")
}
