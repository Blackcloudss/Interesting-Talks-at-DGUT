package logic

import (
	"context"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/repo"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
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

type FollowLogic struct {
}

func NewFollowLogic() *FollowLogic {
	return &FollowLogic{}
}

// Follow 关注用户
func (l *FollowLogic) Follow(ctx context.Context, req types.FollowReq) error {
	tx := global.DB.Begin()
	if tx.Error != nil {
		zlog.CtxErrorf(ctx, "Failed to start transaction: %v", tx.Error)
		return response.ErrResp(tx.Error, codeFollowFailed)
	}

	// 检查是否已经关注
	if repo.IsFollowing(tx, req.FollowerID, req.FollowedID) {
		tx.Rollback()
		zlog.CtxInfof(ctx, "Follow skipped: already following (followerID: %d, followedID: %d)", req.FollowerID, req.FollowedID)
		return response.ErrResp(nil, codeAlreadyFollowing)
	}

	// 插入关注关系
	if err := repo.Follow(tx, req.FollowerID, req.FollowedID); err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "Follow failed: %v", err)
		return response.ErrResp(err, codeFollowFailed)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "Failed to commit transaction: %v", err)
		return response.ErrResp(err, codeFollowFailed)
	}

	zlog.CtxInfof(ctx, "Followed successfully (followerID: %d, followedID: %d)", req.FollowerID, req.FollowedID)
	return nil
}

// Unfollow 取消关注
func (l *FollowLogic) Unfollow(ctx context.Context, req types.UnfollowReq) error {
	tx := global.DB.Begin()
	if tx.Error != nil {
		zlog.CtxErrorf(ctx, "Failed to start transaction: %v", tx.Error)
		return response.ErrResp(tx.Error, codeUnfollowFailed)
	}

	// 检查是否已经取消关注
	if !repo.IsFollowing(tx, req.FollowerID, req.FollowedID) {
		tx.Rollback()
		zlog.CtxInfof(ctx, "Unfollow skipped: not following (followerID: %d, followedID: %d)", req.FollowerID, req.FollowedID)
		return response.ErrResp(nil, codeNotFollowing)
	}

	// 删除关注关系
	if err := repo.Unfollow(tx, req.FollowerID, req.FollowedID); err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "Unfollow failed: %v", err)
		return response.ErrResp(err, codeUnfollowFailed)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "Failed to commit transaction: %v", err)
		return response.ErrResp(err, codeUnfollowFailed)
	}

	zlog.CtxInfof(ctx, "Unfollowed successfully (followerID: %d, followedID: %d)", req.FollowerID, req.FollowedID)
	return nil
}

// GetFollowings 获取用户关注的用户列表
func (l *FollowLogic) GetFollowings(ctx context.Context, req types.GetFollowingsReq) (resp types.GetFollowingsResp, err error) {
	users, err := repo.GetFollowings(req.UserID)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get followings for user (userID: %d): %v", req.UserID, err)
		return types.GetFollowingsResp{}, response.ErrResp(err, codeGetFollowingsFailed)
	}

	zlog.CtxInfof(ctx, "Followings retrieved successfully (userID: %d)", req.UserID)
	return types.GetFollowingsResp{List: users}, nil
}

// GetFollowers 获取用户的粉丝列表
func (l *FollowLogic) GetFollowers(ctx context.Context, req types.GetFollowersReq) (resp types.GetFollowersResp, err error) {
	users, err := repo.GetFollowers(req.UserID)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get followers for user (userID: %d): %v", req.UserID, err)
		return types.GetFollowersResp{}, response.ErrResp(err, codeGetFollowersFailed)
	}

	zlog.CtxInfof(ctx, "Followers retrieved successfully (userID: %d)", req.UserID)
	return types.GetFollowersResp{List: users}, nil
}
