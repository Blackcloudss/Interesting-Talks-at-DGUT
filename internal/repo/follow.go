package repo

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"gorm.io/gorm"
)

// Follow 关注用户
func Follow(tx *gorm.DB, followerID, followedID uint64) error {
	return tx.Create(&model.Follow{
		FollowerID: followerID,
		FollowedID: followedID,
	}).Error
}

// Unfollow 取消关注
func Unfollow(tx *gorm.DB, followerID, followedID uint64) error {
	return tx.Where("follower_id = ? AND followed_id = ?", followerID, followedID).Delete(&model.Follow{}).Error
}

// GetFollowings 获取用户关注的用户列表
func GetFollowings(userID uint64) ([]model.User, error) {
	var users []model.User
	err := global.DB.Model(&model.User{}).
		Joins("inner join user_follows on user_follows.followed_id = users.id").
		Where("user_follows.follower_id = ?", userID).
		Find(&users).Error
	return users, err
}

// GetFollowers 获取用户的粉丝列表
func GetFollowers(userID uint64) ([]model.User, error) {
	var users []model.User
	err := global.DB.Model(&model.User{}).
		Joins("inner join user_follows on user_follows.follower_id = users.id").
		Where("user_follows.followed_id = ?", userID).
		Find(&users).Error
	return users, err
}

// IsFollowing 检查是否已经关注
func IsFollowing(followerID, followedID uint64) bool {
	var count int64
	global.DB.Model(&model.Follow{}).
		Where("follower_id = ? AND followed_id = ?", followerID, followedID).
		Count(&count)
	return count > 0
}
