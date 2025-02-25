package logic

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"gorm.io/gorm"
)

// @Title        hook.go
// @Description
// @Create       XdpCs 2025-02-16 上午1:32
// @Update       XdpCs 2025-02-16 上午1:32
// RegisterHook 注册 GORM 钩子
func RegisterHook(db *gorm.DB) {
	zlog.Infof("注册 GORM hooks...")
	db.Callback().Create().Before("gorm:Create").Register("before_create_BaseModel", BeforeCreateBaseModel)
}

func BeforeCreateBaseModel(db *gorm.DB) {
	if db.Statement.Schema != nil {
		if baseModel, ok := db.Statement.Model.(*model.CommonModel); ok {
			baseModel.BeforeCreate(db)
		}
	}
}
