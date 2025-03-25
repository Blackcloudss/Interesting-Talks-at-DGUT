package repo

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"gorm.io/gorm"
)

type FollowRepo struct {
	DB *gorm.DB
}

func NewFollowRepo(db *gorm.DB) *FollowRepo {
	return &FollowRepo{
		DB: db,
	}
}

// Follow 关注用户
func (r *FollowRepo) Follow(followerID, followedID int64) error {
	// 检查是否已经关注
	if r.IsFollowing(followerID, followedID) {
		return nil // 如果已经关注，直接返回
	}

	// 检查对方是否已经关注了当前用户
	isFollowedBack := r.IsFollowing(followedID, followerID)

	// 创建关注记录
	follow := &model.Follow{
		FollowerID: followerID,
		FollowedID: followedID,
		IsFriend:   isFollowedBack, // 如果对方已经关注了当前用户，标记为好友
	}
	if err := r.DB.Create(follow).Error; err != nil {
		return err
	}

	// 如果对方已经关注了当前用户，更新对方的关注记录，将对方的关注关系也标记为好友
	if isFollowedBack {
		return r.DB.Model(&model.Follow{}).
			Where("follower_id = ? AND followed_id = ?", followedID, followerID).
			Update("is_friend", true).Error
	}

	return nil
}

// Unfollow 取消关注
func (r *FollowRepo) Unfollow(followerID, followedID int64) error {
	return r.DB.Where("follower_id = ? AND followed_id = ?", followerID, followedID).Delete(&model.Follow{}).Error
}

// GetFollowings 获取用户关注的用户列表
func (r *FollowRepo) GetFollowings(userID int64) ([]types.FollowInfo, error) {
	var users []types.FollowInfo
	err := r.DB.Model(&model.Follow{}).
		Select("user_display.id, user_display.nickname, user_display.avatar, follow.created_at as followed_at").
		Joins("INNER JOIN user_display ON user_display.id = follow.followed_id").
		Where("follower_id = ?", userID).
		Scan(&users).Error
	return users, err
}

// GetFollowers 获取用户的粉丝列表
func (r *FollowRepo) GetFollowers(userID int64) ([]types.FollowInfo, error) {
	var users []types.FollowInfo
	err := r.DB.Model(&model.Follow{}).
		Select("user_display.id, user_display.nickname, user_display.avatar, follow.created_at as followed_at").
		Joins("INNER JOIN user_display ON user_display.id = follow.follower_id").
		Where("followed_id = ?", userID).
		Scan(&users).Error
	return users, err
}

// IsFollowing 检查是否已经关注
func (r *FollowRepo) IsFollowing(followerID, followedID int64) bool {
	var count int64
	r.DB.Model(&model.Follow{}).
		Where("follower_id = ? AND followed_id = ?", followerID, followedID).
		Count(&count)
	return count > 0
}
