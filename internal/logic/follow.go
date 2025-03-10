package logic

import (
	"context"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/repo"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
)

// 定义内部逻辑错误
var (
	codeFollowFailed        = response.MsgCode{Code: 40035, Msg: "关注失败"}
	codeUnfollowFailed      = response.MsgCode{Code: 40036, Msg: "取消关注失败"}
	codeAlreadyFollowing    = response.MsgCode{Code: 40037, Msg: "已经关注该用户"}
	codeNotFollowing        = response.MsgCode{Code: 40038, Msg: "未关注该用户"}
	codeGetFollowingsFailed = response.MsgCode{Code: 40039, Msg: "获取关注列表失败"}
	codeGetFollowersFailed  = response.MsgCode{Code: 40040, Msg: "获取粉丝列表失败"}
)

// Follow 关注用户
func Follow(ctx context.Context, followerID, followedID uint64) error {
	tx := global.DB.Begin()
	if tx.Error != nil {
		zlog.CtxErrorf(ctx, "Failed to start transaction: %v", tx.Error)
		return response.ErrResp(tx.Error, codeFollowFailed)
	}

	// 检查是否已经关注
	if repo.IsFollowing(followerID, followedID) {
		tx.Rollback()
		zlog.CtxInfof(ctx, "Follow skipped: already following (followerID: %d, followedID: %d)", followerID, followedID)
		return response.ErrResp(nil, codeAlreadyFollowing)
	}

	// 插入关注关系
	if err := repo.Follow(tx, followerID, followedID); err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "Follow failed: %v", err)
		return response.ErrResp(err, codeFollowFailed)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "Failed to commit transaction: %v", err)
		return response.ErrResp(err, codeFollowFailed)
	}

	zlog.CtxInfof(ctx, "Followed successfully (followerID: %d, followedID: %d)", followerID, followedID)
	return nil
}

// Unfollow 取消关注
func Unfollow(ctx context.Context, followerID, followedID uint64) error {
	tx := global.DB.Begin()
	if tx.Error != nil {
		zlog.CtxErrorf(ctx, "Failed to start transaction: %v", tx.Error)
		return response.ErrResp(tx.Error, codeUnfollowFailed)
	}

	// 检查是否已经取消关注
	if !repo.IsFollowing(followerID, followedID) {
		tx.Rollback()
		zlog.CtxInfof(ctx, "Unfollow skipped: not following (followerID: %d, followedID: %d)", followerID, followedID)
		return response.ErrResp(nil, codeNotFollowing)
	}

	// 删除关注关系
	if err := repo.Unfollow(tx, followerID, followedID); err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "Unfollow failed: %v", err)
		return response.ErrResp(err, codeUnfollowFailed)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "Failed to commit transaction: %v", err)
		return response.ErrResp(err, codeUnfollowFailed)
	}

	zlog.CtxInfof(ctx, "Unfollowed successfully (followerID: %d, followedID: %d)", followerID, followedID)
	return nil
}

// GetFollowings 获取用户关注的用户列表
func GetFollowings(ctx context.Context, userID uint64) ([]model.User, error) {
	users, err := repo.GetFollowings(userID)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get followings for user (userID: %d): %v", userID, err)
		return nil, response.ErrResp(err, codeGetFollowingsFailed)
	}

	zlog.CtxInfof(ctx, "Followings retrieved successfully (userID: %d)", userID)
	return users, nil
}

// GetFollowers 获取用户的粉丝列表
func GetFollowers(ctx context.Context, userID uint64) ([]model.User, error) {
	users, err := repo.GetFollowers(userID)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get followers for user (userID: %d): %v", userID, err)
		return nil, response.ErrResp(err, codeGetFollowersFailed)
	}

	zlog.CtxInfof(ctx, "Followers retrieved successfully (userID: %d)", userID)
	return users, nil
}
