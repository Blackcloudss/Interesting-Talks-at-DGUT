package repo

import (
	"errors"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils"
	"gorm.io/gorm"
	"time"
)

const ROLE = "role"

// @Title        casbin.go
// @Description
// @Create       XdpCs 2025-03-11 上午1:15
// @Update       XdpCs 2025-03-11 上午1:15

type CasbinRepo struct {
	DB *gorm.DB
}

func NewCasbinRepo(db *gorm.DB) *CasbinRepo {
	return &CasbinRepo{
		DB: db,
	}
}

// CheckUserPermission
//
//	@Description:
//	@receiver r
//	@param Url
//	@param UserId
//	@return bool
//	@return error
func (r *CasbinRepo) CheckUserPermission(Durl string, UserId int64) (IsExist bool, err error) {
	defer utils.RecordTime(time.Now())()
	//查找该用户所扮演的角色
	var Role string
	err = r.DB.Model(&model.UserMsg{}).
		Select(ROLE).
		Where(&model.UserMsg{
			CommonModel: model.CommonModel{
				ID: UserId,
			},
		}).
		First(&Role).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		zlog.Errorf("查询用户对应的角色失败：%v", err)
		return false, err
	}
	if Role == global.TOURIST {
		IsExist, err = r.Contains(global.TOURIST_URLS, Durl)
		if err != nil {
			zlog.Errorf("查询失败：: %v", err)
			return false, err
		}
		return IsExist, nil
	} else if Role == global.STUDENT {
		IsExist, err = r.Contains(global.STUDENT_URLS, Durl)
		if err != nil {
			zlog.Errorf("查询失败：: %v", err)
			return false, err
		}
		return IsExist, nil
	} else if Role == global.MANAGER {
		IsExist, err = r.Contains(global.MANAGER_URLS, Durl)
		if err != nil {
			zlog.Errorf("查询失败：: %v", err)
			return false, err
		}
		return IsExist, nil
	} else {
		zlog.Errorf("该用户所具有的角色不存在：: %v", nil)
		return false, nil
	}
}

// 判断 当前 url 是否在 该角色所拥有的urls里
func (r *CasbinRepo) Contains(Urls []string, Durl string) (bool, error) {
	for _, Url := range Urls {
		if Url == Durl {
			return true, nil
		}
	}
	return false, nil
}
