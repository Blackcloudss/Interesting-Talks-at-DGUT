package repo

import (
	"fmt"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	IS_FOLLOWING = "is_following"
	IS_FRIEND    = "is_friend"
)

type FollowRepo struct {
	DB *gorm.DB
}

func NewFollowRepo(db *gorm.DB) *FollowRepo {
	return &FollowRepo{
		DB: db,
	}
}

func (r *FollowRepo) Follow(followerID, followedID int64) (*types.FollowResp, error) {
	var isFollowing bool

	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查询关注状态
	err := tx.Model(&model.Follow{}).
		Select(IS_FOLLOWING).
		Where(fmt.Sprintf("%s = ? AND %s = ?", FOLLOWER_ID, FOLLOWED_ID), followerID, followedID).
		First(&isFollowing).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 初始化关注记录
			isFollowing = true
			follow := model.Follow{
				FollowerID:  followerID,
				FollowedID:  followedID,
				IsFriend:    r.IsFollowing(followedID, followerID),
				IsFollowing: true,
			}
			if err := tx.Create(&follow).Error; err != nil {
				tx.Rollback()
				zlog.Errorf("创建关注关系失败: %v", err)
				return nil, err
			}

			// 如果对方已经关注了当前用户，更新对方的关注记录为好友
			if follow.IsFriend {
				if err := tx.Model(&model.Follow{}).
					Where(fmt.Sprintf("%s = ? AND %s = ?", FOLLOWER_ID, FOLLOWED_ID), followedID, followerID).
					Update(IS_FRIEND, true).Error; err != nil {
					tx.Rollback()
					zlog.Errorf("更新对方关注关系为好友失败: %v", err)
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
		isFollowing = !isFollowing
		if isFollowing {
			// 重新关注
			follow := model.Follow{
				FollowerID:  followerID,
				FollowedID:  followedID,
				IsFriend:    r.IsFollowing(followedID, followerID),
				IsFollowing: true,
			}
			if err := tx.Create(&follow).Error; err != nil {
				tx.Rollback()
				zlog.Errorf("重新关注失败: %v", err)
				return nil, err
			}
		} else {
			// 取消关注
			if err := tx.Where(fmt.Sprintf("%s = ? AND %s = ?", FOLLOWER_ID, FOLLOWED_ID), followerID, followedID).
				Delete(&model.Follow{}).Error; err != nil {
				tx.Rollback()
				zlog.Errorf("取消关注失败: %v", err)
				return nil, err
			}
		}

		// 更新好友关系状态
		if err := tx.Model(&model.Follow{}).
			Where(fmt.Sprintf("%s = ? AND %s = ?", FOLLOWER_ID, FOLLOWED_ID), followedID, followerID).
			Update(IS_FRIEND, isFollowing && r.IsFollowing(followedID, followerID)).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			zlog.Errorf("更新好友关系状态失败: %v", err)
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &types.FollowResp{
		IsFollowing: isFollowing,
	}, nil
}
func (r *FollowRepo) IsFollowing(followerID, followedID int64) bool {
	var count int64
	r.DB.Model(&model.Follow{}).
		Where(fmt.Sprintf("%s = ? AND %s = ? AND %s = ?", FOLLOWER_ID, FOLLOWED_ID, IS_FOLLOWING), followerID, followedID, true).
		Count(&count)
	return count > 0
}

// 关注列表
func (r *FollowRepo) GetFollowings(userID int64) ([]types.FollowInfo, error) {
	var users []types.FollowInfo
	err := r.DB.Model(&model.Follow{}).
		Select("user_display.id as user_id, user_display.nickname, user_display.avatar, follow.created_at as followed_at").
		Joins("JOIN user_display ON user_display.id = follow.followed_id").
		Where(fmt.Sprintf("%s = ? AND %s = ?", FOLLOWER_ID, IS_FOLLOWING), userID, true).
		Scan(&users).Error
	return users, err
}

// 粉丝列表
func (r *FollowRepo) GetFollowers(userID int64) ([]types.FollowInfo, error) {
	var users []types.FollowInfo
	err := r.DB.Model(&model.Follow{}).
		Select("user_display.id as user_id, user_display.nickname, user_display.avatar, follow.created_at as followed_at").
		Joins("JOIN user_display ON user_display.id = follow.follower_id").
		Where(fmt.Sprintf("%s = ? AND %s = ?", FOLLOWED_ID, IS_FOLLOWING), userID, true).
		Scan(&users).Error
	return users, err
}
