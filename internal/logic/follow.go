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

type FollowLogic struct{}

func NewFollowLogic() *FollowLogic {
	return &FollowLogic{}
}

// Follow 关注用户
func (l *FollowLogic) Follow(ctx context.Context, req types.FollowReq, UserID int64) (*types.FollowResp, error) {
	// 检查是否已经关注
	if repo.IsFollowing(global.DB, UserID, req.FollowedID) {
		zlog.CtxInfof(ctx, "Follow skipped: already following (followerID: %d, followedID: %d)", UserID, req.FollowedID)
		return nil, response.ErrResp(nil, codeAlreadyFollowing)
	}

	// 插入关注关系
	if err := repo.Follow(global.DB, UserID, req.FollowedID); err != nil {
		zlog.CtxErrorf(ctx, "Follow failed: %v", err)
		return nil, response.ErrResp(err, codeFollowFailed)
	}

	zlog.CtxInfof(ctx, "Followed successfully (followerID: %d, followedID: %d)", UserID, req.FollowedID)
	return &types.FollowResp{Success: true}, nil
}

// Unfollow 取消关注
func (l *FollowLogic) Unfollow(ctx context.Context, req types.UnfollowReq, UserID int64) (*types.UnfollowResp, error) {
	// 检查是否已经取消关注
	if !repo.IsFollowing(global.DB, UserID, req.FollowedID) {
		zlog.CtxInfof(ctx, "Unfollow skipped: not following (followerID: %d, followedID: %d)", UserID, req.FollowedID)
		return nil, response.ErrResp(nil, codeNotFollowing)
	}

	// 删除关注关系
	if err := repo.Unfollow(global.DB, UserID, req.FollowedID); err != nil {
		zlog.CtxErrorf(ctx, "Unfollow failed: %v", err)
		return nil, response.ErrResp(err, codeUnfollowFailed)
	}

	zlog.CtxInfof(ctx, "Unfollowed successfully (followerID: %d, followedID: %d)", UserID, req.FollowedID)
	return &types.UnfollowResp{Success: true}, nil
}

// GetFollowings 获取用户关注的用户列表
func (l *FollowLogic) GetFollowings(ctx context.Context, req types.GetFollowingsReq, UserID int64) (*types.GetFollowingsResp, error) {
	users, err := repo.GetFollowings(UserID)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get followings for user (userID: %d): %v", UserID, err)
		return nil, response.ErrResp(err, codeGetFollowingsFailed)
	}

	zlog.CtxInfof(ctx, "Followings retrieved successfully (userID: %d)", UserID)
	return &types.GetFollowingsResp{Followers: users}, nil
}

// GetFollowers 获取用户的粉丝列表
func (l *FollowLogic) GetFollowers(ctx context.Context, req types.GetFollowersReq, UserID int64) (*types.GetFollowersResp, error) {
	users, err := repo.GetFollowers(UserID)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get followers for user (userID: %d): %v", UserID, err)
		return nil, response.ErrResp(err, codeGetFollowersFailed)
	}

	zlog.CtxInfof(ctx, "Followers retrieved successfully (userID: %d)", UserID)
	return &types.GetFollowersResp{Fans: users}, nil
}
