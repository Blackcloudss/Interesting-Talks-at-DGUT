package repo

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
)

// GetUserRole 获取用户角色
func GetUserRole(userID uint64) (string, error) {
	var user model.User
	if err := global.DB.First(&user, userID).Error; err != nil {
		return "", err
	}
	return user.Role, nil
}
