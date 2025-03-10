package repo

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"gorm.io/gorm"
)

// @Title        login.go
// @Description
// @Create       XdpCs 2025-03-10 上午1:22
// @Update       XdpCs 2025-03-10 上午1:22
type WechatLoginRepo struct {
	DB *gorm.DB
}

func NewWechatLoginRepo(db *gorm.DB) *WechatLoginRepo {
	return &WechatLoginRepo{
		DB: db,
	}
}

func (r WechatLoginRepo) GetUserId(userId string) (testUser model.Test, err error) {

	return
}
