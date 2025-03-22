package repo

import (
	"errors"
	"fmt"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/snowflake"
	"gorm.io/gorm"
)

const (
	USER_DISPLAY = "user_display"
	USER_COMMON  = "user_common"
	NICKNAME     = "nickname"
	AVATAR       = "avatar"
	TAG          = "tag"
	SEX          = "sex"
	BIRTHDAY     = "birthday"
	SIGN         = "sign"
	NAME         = "name"
	STUDENT_ID   = "student_id"
	ACADEMY      = "academy"
	GRADE        = "grade"
	MAJOR        = "major"
	PHONE        = "phone"
	USERDISPLAY  = "UserDisplay"
)

// @Title        middle.go
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

			//通过雪花算法生成随机昵称id
			NickId := snowflake.GetInt12Id(global.Node)
			// 创建用户展示表
			UserDisplay := &model.UserDisplay{
				OpenId:   Openid,
				Nickname: fmt.Sprintf("莞星人%d", NickId),
			}
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
	return nil
}

// GetCommonProfile
//
//	@Description: 查找用户基本信息
//	@receiver r
//	@param UserId
//	@return UserCommon
//	@return err
func (r *UserRepo) GetCommonProfile(UserId int64) (resp types.GetCommonProfileResp, err error) {
	err = r.DB.Preload(USERDISPLAY).
		Select(fmt.Sprintf("%s.%s,%s.%s,%s.%s,%s.%s,%s.%s,%s.%s",
			USER_DISPLAY, NICKNAME, USER_DISPLAY, AVATAR, USER_DISPLAY, TAG,
			USER_COMMON, BIRTHDAY, USER_COMMON, SEX, USER_COMMON, SIGN,
		)).
		Where(fmt.Sprintf("%s = ?", USER_ID), UserId).
		First(&resp).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.Errorf("用户不存在：%v", err)
			return resp, err
		}
		zlog.Errorf("查询用户基本信息失败：%v", err)
		return resp, err
	}
	return resp, nil
}

// UpdateCommonProfile
//
//	@Description: 更新用户隐私信息
//	@receiver r
//	@param UserId
//	@param req
//	@return err
func (r *UserRepo) UpdateCommonProfile(UserId int64, Avatar string, req types.UpdateCommonProfileReq) (err error) {
	err = r.DB.Model(&model.UserDisplay{}).
		Where(fmt.Sprintf("%s = ?", ID), UserId).
		Updates(&model.UserDisplay{
			Avatar:   Avatar,
			Nickname: req.Nickname,
		}).Error
	if err != nil {
		zlog.Errorf("更新用户基本信息失败：%v", err)
		return err
	}

	err = r.DB.Model(&model.UserCommon{}).
		Where(fmt.Sprintf("%s = ?", USER_ID), UserId).
		Updates(&model.UserCommon{
			Birthday: req.Birthday,
			Sex:      req.Sex,
			Sign:     req.Sign,
		}).Error
	if err != nil {
		zlog.Errorf("更新用户基本信息失败：%v", err)
		return err
	}
	return nil
}

// GetUserPrivateInfo
//
//	@Description: 查找用户隐私信息
//	@receiver r
//	@param UserId
//	@return UserPrivate
//	@return err
func (r *UserRepo) GetPrivateProfile(UserId int64) (resp types.GetPrivateProfileResp, err error) {
	err = r.DB.Model(&model.UserPrivate{}).
		Select(fmt.Sprintf("%s, %s, %s, %s, %s, %s", NAME, STUDENT_ID, ACADEMY, GRADE, MAJOR, PHONE)).
		Where(&model.UserPrivate{
			UserID: UserId,
		}).
		First(&resp).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.Errorf("用户不存在：%v", err)
			return resp, err
		}
		zlog.Errorf("查询用户隐私信息失败：%v", err)
		return resp, err
	}
	return resp, nil
}

// UpdatePrivateProfile
//
//	@Description: 更改用户隐私信息
//	@receiver r
//	@param UserId
//	@param req
//	@return err
func (r *UserRepo) UpdatePrivateProfile(UserId int64, req types.UpdatePrivateProfileReq) (Role string, err error) {
	// 更改用户隐私信息
	err = r.DB.Model(&model.UserPrivate{}).
		Where(&model.UserPrivate{
			UserID: UserId,
		}).
		Updates(&model.UserPrivate{
			Name:      req.Name,
			StudentId: req.StudentId,
			Academy:   req.Academy,
			Grade:     req.Grade,
			Major:     req.Major,
			Phone:     req.Phone,
		}).Error
	if err != nil {
		zlog.Errorf("更新用户隐私信息失败：%v", err)
		return "", err
	}
	//更改用户身份
	err = r.DB.Model(&model.UserDisplay{}).
		Where(&model.UserDisplay{
			CommonModel: model.CommonModel{
				ID: UserId,
			},
		}).
		Updates(&model.UserDisplay{
			Role: STUDENT,
		}).Error
	if err != nil {
		zlog.Errorf("更新用户身份失败：%v", err)
		return "", err
	}
	Role = STUDENT

	return Role, nil
}

// GetUserRole
//
//	@Description: 获取用户角色身份
//	@receiver r
//	@param userID
//	@return string
//	@return error
func (r *UserRepo) GetUserRole(userID int64) (role string, err error) {
	err = global.DB.Model(&model.UserDisplay{}).
		Select(ROLE).
		Where(model.UserDisplay{
			CommonModel: model.CommonModel{
				ID: userID,
			},
		}).
		First(&role).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.Errorf("用户不存在：%v", err)
			return "", err
		}
		zlog.Errorf("查询用户角色失败：%v", err)
		return "", err
	}
	return role, nil
}

// UpdateOtherRole
//
//	@Description: 更改用户角色身份
//	@receiver r
//	@param UserID
//	@param Role
//	@return error
func (r *UserRepo) UpdateOtherRole(UserID int64, Role string) error {
	err := global.DB.Model(&model.UserDisplay{}).
		Where(model.UserDisplay{
			CommonModel: model.CommonModel{
				ID: UserID,
			},
		}).
		Updates(&model.UserDisplay{
			Role: Role,
		}).Error
	if err != nil {
		zlog.Errorf("更新用户角色失败：%v", err)
		return err
	}
	return nil
}
