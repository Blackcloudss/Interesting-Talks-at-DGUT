package repo

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"gorm.io/gorm"
)

// @Title        phone.go
// @Description
// @Create       XdpCs 2025-03-17 下午7:11
// @Update       XdpCs 2025-03-17 下午7:11
type PhoneRepo struct {
	DB *gorm.DB
}

func NewPhoneRepo(db *gorm.DB) *PhoneRepo {
	return &PhoneRepo{
		DB: db,
	}
}

func (r *PhoneRepo) SavePhone(userid int64, phone string) (err error) {
	if err = r.DB.Model(&model.UserDisplay{}).
		Where("user_id = ?", userid).
		Update("phone", phone).
		Error; err != nil {
		return
	}
	return
}
