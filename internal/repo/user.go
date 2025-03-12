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

// @Title        tz_user.go
// @Description
// @Create       XdpCs 2025-03-10 下午1:59
// @Update       XdpCs 2025-03-10 下午1:59

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}

// 判断用户是否存在
func (r *UserRepo) JudgeUser(Openid string) (int64, error) {

	var UserID int64

	err := r.DB.Model(&model.UserDisplay{}).
		Select(Id).
		Where(&model.UserDisplay{
			OpenId: Openid,
		}).
		First(&UserID).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 开启事务
			tx := r.DB.Begin()
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()

			// 创建用户展示表
			userDisplay := &model.UserDisplay{OpenId: Openid}
			if err = tx.Create(userDisplay).Error; err != nil {
				tx.Rollback()
				zlog.Errorf("创建用户展示表失败：%v", err)
				return 0, err
			}

			// 使用事务创建后续表
			tables := []interface{}{
				&model.UserCommon{OpenId: Openid},
				&model.UserPrivate{OpenId: Openid},
				&model.UserAuth{OpenId: Openid},
			}

			for _, table := range tables {
				if err = tx.Create(table).Error; err != nil {
					tx.Rollback()
					zlog.Errorf("创建用户表失败：%v", err)
					return 0, err
				}
			}

			// 提交事务
			if err = tx.Commit().Error; err != nil {
				zlog.Errorf("事务提交失败：%v", err)
				return 0, err
			}

			// 直接使用创建后获得的ID
			return userDisplay.ID, nil
		}
		return 0, err
	}
	return UserID, nil
}
