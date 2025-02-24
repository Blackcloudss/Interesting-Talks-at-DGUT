package model

import (
	"gorm.io/gorm"
	"time"
)

// @Title        common.go
// @Description
// @Create       XdpCs 2025-02-24 下午11:00
// @Update       XdpCs 2025-02-24 下午11:00
// CommonModel 每张表都有的四个东西，最好不要用 gorm.model（虽然他们一模一样）
type CommonModel struct {
	ID        int64 `gorm:"primaryKey;column:id;type:bigint"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (b *CommonModel) BeforeCreate(db *gorm.DB) error {
	// 生成雪花ID
	if b.ID == 0 {
		b.ID = snowflake.GetInt12Id(global.Node)
	}
	return nil
}
