package repo

import (
	"errors"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"gorm.io/gorm"
)

const (
	Id = "id"
)

// @Title        user.go
// @Description
// @Create       XdpCs 2025-03-10 下午1:59
// @Update       XdpCs 2025-03-10 下午1:59

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) JudgeUser(Openid string) (int64, error) {

	var UserID int64

	err := r.DB.Model(&model.UserMsg{}).
		Select(Id).
		Where(&model.UserMsg{
			OpenId: Openid,
		}).
		First(&UserID).
		Error
	if err != nil {
		// 如果数据库中没有该记录，则创建新用户 存放 OpenId
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = r.DB.Model(&model.UserMsg{}).
				Create(&model.UserMsg{
					OpenId: Openid,
				}).Error
			if err != nil {
				zlog.Errorf("创建用户失败：: %v", err)
				return 0, err
			}
			// 查询新创建的用户的 ID
			err = r.DB.Model(&model.UserMsg{}).
				Select(Id).
				Where(&model.UserMsg{
					OpenId: Openid,
				}).
				First(&UserID).
				Error
			return UserID, nil
		}
		zlog.Errorf("查询失败：: %v", err)
		return 0, err
	}
	return UserID, nil
}
