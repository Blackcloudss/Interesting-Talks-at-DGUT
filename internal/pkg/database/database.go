package database

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/configs"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"gorm.io/gorm"
)

// @Title        database.go
// @Description
// @Create       XdpCs 2025-02-24 下午11:24
// @Update       XdpCs 2025-02-24 下午11:24
type DataBase interface {
	GetDsn(config configs.Config) string
	InitDataBase(config configs.Config) (*gorm.DB, error)
}

func InitDataBases(base DataBase, config configs.Config) {
	var err error
	global.DB, err = base.InitDataBase(config)
	if err != nil {
		zlog.Fatalf("无法初始化数据库 %v", err)
		return
	}
	zlog.Infof("初始化数据库成功！")
	//对该数据库注册 hook
	logic.RegisterHook(global.DB)
	return
}
