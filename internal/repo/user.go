package repo

import (
	"errors"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"gorm.io/gorm"
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
		Select(ID).
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
			UserDisplay := &model.UserDisplay{OpenId: Openid}
			if err = tx.Create(UserDisplay).Error; err != nil {
				tx.Rollback()
				zlog.Errorf("创建用户展示表失败：%v", err)
				return 0, err
			}
			zlog.Infof("创建用户展示表成功，ID：%d", UserDisplay.ID)
			UserID = UserDisplay.ID

			// 使用事务创建后续表
			tables := []interface{}{
				&model.UserCommon{UserID: UserID},
				&model.UserPrivate{UserID: UserID},
				&model.UserAuth{UserID: UserID},
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
			zlog.Infof("创建用户表成功，ID：%d", UserDisplay.ID)
			return UserID, nil
		}
		return 0, err
	}
	return UserID, nil
}

// GetOpenId
//
//	@Description: 根据UserID获取用户拥有的OpenId
//	@receiver r
//	@param UserId
//	@return OpenId
//	@return Err
func (r *UserRepo) GetOpenId(UserId int64) (OpenId int64, err error) {
	err = r.DB.Model(&model.UserDisplay{}).
		Select(OPEN_ID).
		Where(&model.UserDisplay{
			CommonModel: model.CommonModel{
				ID: UserId,
			},
		}).
		First(&OpenId).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.Errorf("用户不存在：%v", err)
			return 0, err
		}
		zlog.Errorf("查询用户对应的OpenId失败：%v", err)
		return 0, err
	}
	return
}

// SaveUserInfo
//
//	@Description: 保存用户信息到数据库中
//	@receiver r
//	@param UserId
//	@param User
//	@return err
func (r *UserRepo) SaveUserInfo(UserId int64, User types.UserInfo) (err error) {
	err = r.DB.Model(&model.UserDisplay{}).
		Where(&model.UserDisplay{
			CommonModel: model.CommonModel{
				ID: UserId,
			},
		}).
		Updates(&model.UserDisplay{
			Avatar:   User.Avatar,
			Nickname: User.Nickname,
		}).Error
	if err != nil {
		zlog.Errorf("更新用户信息失败：%v", err)
		return err
	}
	return
}
