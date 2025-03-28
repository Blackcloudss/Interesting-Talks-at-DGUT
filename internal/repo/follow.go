package repo

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	FollowerID  = "follower_id"
	FollowedID  = "followed_id"
	IsFollowing = "is_following"
	IsFriend    = "is_friend"
)

type FollowRepo struct {
	DB *gorm.DB
}

func NewFollowRepo(db *gorm.DB) *FollowRepo {
	return &FollowRepo{DB: db}
}

func (r *FollowRepo) Follow(followerID, followedID int64) (*types.FollowResp, error) {
	if followerID == followedID {
		return nil, errors.New("不能关注自己")
	}
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var follow model.Follow
	err := tx.Where("follower_id = ? AND followed_id = ?", followerID, followedID).
		First(&follow).Error

	var isFollowing bool
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新的关注关系
			isFollowing = true
			follow = model.Follow{
				FollowerID:  followerID,
				FollowedID:  followedID,
				IsFollowing: true,
				IsFriend:    r.IsFollowing(followedID, followerID),
			}
			if err := tx.Create(&follow).Error; err != nil {
				tx.Rollback()
				zlog.Errorf("创建关注关系失败: %v", err)
				return nil, err
			}

			// 如果对方已经关注了当前用户，更新对方记录为好友
			if follow.IsFriend {
				if err := tx.Model(&model.Follow{}).
					Where("follower_id = ? AND followed_id = ?", followedID, followerID).
					Update("is_friend", true).Error; err != nil {
					tx.Rollback()
					zlog.Errorf("更新对方好友状态失败: %v", err)
					return nil, err
				}
			}
		} else {
			tx.Rollback()
			zlog.Errorf("查询关注状态失败: %v", err)
			return nil, err
		}
	} else {
		// 切换关注状态
		isFollowing = !follow.IsFollowing
		if isFollowing {
			// 重新关注
			if err := tx.Model(&follow).
				Updates(map[string]interface{}{
					"is_following": true,
					"is_friend":    r.IsFollowing(followedID, followerID),
				}).Error; err != nil {
				tx.Rollback()
				zlog.Errorf("更新关注状态失败: %v", err)
				return nil, err
			}
		} else {
			// 取消关注
			if err := tx.Model(&follow).
				Update("is_following", false).Error; err != nil {
				tx.Rollback()
				zlog.Errorf("取消关注失败: %v", err)
				return nil, err
			}
		}

		// 更新对方的好友状态
		if err := tx.Model(&model.Follow{}).
			Where("follower_id = ? AND followed_id = ?", followedID, followerID).
			Update("is_friend", false).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			zlog.Errorf("更新对方好友状态失败: %v", err)
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &types.FollowResp{IsFollowing: isFollowing}, nil
}

func (r *FollowRepo) IsFollowing(followerID, followedID int64) bool {
	var count int64
	r.DB.Model(&model.Follow{}).
		Where("follower_id = ? AND followed_id = ? AND is_following = ?",
			followerID, followedID, true).
		Count(&count)
	return count > 0
}

func (r *FollowRepo) GetFollowings(userID int64) ([]types.FollowInfo, error) {
	var users []types.FollowInfo
	err := r.DB.Model(&model.Follow{}).
		Select("user_display.id as user_id, user_display.nickname, user_display.avatar, follow.created_at as followed_at, follow.is_following").
		Joins("JOIN user_display ON user_display.id = follow.followed_id").
		Where("follow.follower_id = ? AND follow.is_following = ?", userID, true).
		Scan(&users).Error
	return users, err
}

func (r *FollowRepo) GetFollowers(userID int64) ([]types.FollowInfo, error) {
	var users []types.FollowInfo
	err := r.DB.Model(&model.Follow{}).
		Select("user_display.id as user_id, user_display.nickname, user_display.avatar, follow.created_at as followed_at, follow.is_following").
		Joins("JOIN user_display ON user_display.id = follow.follower_id").
		Where("follow.followed_id = ? AND follow.is_following = ?", userID, true).
		Scan(&users).Error
	return users, err
}
