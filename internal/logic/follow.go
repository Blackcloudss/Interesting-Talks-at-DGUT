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
	if repo.NewFollowRepo(global.DB).IsFollowing(UserID, req.FollowedID) {
		zlog.CtxInfof(ctx, "Follow skipped: already following (followerID: %d, followedID: %d)", UserID, req.FollowedID)
		return nil, response.ErrResp(nil, codeAlreadyFollowing)
	}

	// 插入关注关系
	if err := repo.NewFollowRepo(global.DB).Follow(UserID, req.FollowedID); err != nil {
		zlog.CtxErrorf(ctx, "Follow failed: %v", err)
		return nil, response.ErrResp(err, codeFollowFailed)
	}

	zlog.CtxInfof(ctx, "Followed successfully (followerID: %d, followedID: %d)", UserID, req.FollowedID)
	return &types.FollowResp{}, nil
}

// Unfollow 取消关注
func (l *FollowLogic) Unfollow(ctx context.Context, req types.UnfollowReq, UserID int64) (*types.UnfollowResp, error) {
	// 检查是否已经取消关注
	if !repo.NewFollowRepo(global.DB).IsFollowing(UserID, req.FollowedID) {
		zlog.CtxInfof(ctx, "Unfollow skipped: not following (followerID: %d, followedID: %d)", UserID, req.FollowedID)
		return nil, response.ErrResp(nil, codeNotFollowing)
	}

	// 删除关注关系
	if err := repo.NewFollowRepo(global.DB).Unfollow(UserID, req.FollowedID); err != nil {
		zlog.CtxErrorf(ctx, "Unfollow failed: %v", err)
		return nil, response.ErrResp(err, codeUnfollowFailed)
	}

	zlog.CtxInfof(ctx, "Unfollowed successfully (followerID: %d, followedID: %d)", UserID, req.FollowedID)
	return &types.UnfollowResp{}, nil
}

// GetFollowings 获取用户关注的用户列表(我的和其他用户的都可以使用）
func (l *FollowLogic) GetFollowings(ctx context.Context, UserID int64) (*types.GetFollowingsResp, error) {
	followings, err := repo.NewFollowRepo(global.DB).GetFollowings(UserID)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get followings for user (userID: %d): %v", UserID, err)
		return nil, response.ErrResp(err, codeGetFollowingsFailed)
	}

	// 构造响应
	resp := &types.GetFollowingsResp{
		Followings: make([]types.FollowInfo, len(followings)),
	}
	for i, user := range followings {
		resp.Followings[i] = types.FollowInfo{
			UserID:       user.UserID,
			Nickname:     user.Nickname,
			Avatar:       user.Avatar,
			FollowedAt:   user.FollowedAt,
			IsFollowedBy: repo.NewFollowRepo(global.DB).IsFollowing(user.UserID, UserID),
		}
	}

	zlog.CtxInfof(ctx, "Followings retrieved successfully (userID: %d)", UserID)
	return resp, nil
}

// GetFollowers 获取用户的粉丝列表(我的和其他用户的都可以使用）
func (l *FollowLogic) GetFollowers(ctx context.Context, UserID int64) (*types.GetFollowersResp, error) {
	followers, err := repo.NewFollowRepo(global.DB).GetFollowers(UserID)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get followers for user (userID: %d): %v", UserID, err)
		return nil, response.ErrResp(err, codeGetFollowersFailed)
	}

	// 构造响应
	resp := &types.GetFollowersResp{
		Followers: make([]types.FollowInfo, len(followers)),
	}
	for i, user := range followers {
		resp.Followers[i] = types.FollowInfo{
			UserID:       user.UserID,
			Nickname:     user.Nickname,
			Avatar:       user.Avatar,
			FollowedAt:   user.FollowedAt,
			IsFollowedBy: repo.NewFollowRepo(global.DB).IsFollowing(user.UserID, UserID),
		}
	}

	zlog.CtxInfof(ctx, "Followers retrieved successfully (userID: %d)", UserID)
	return resp, nil
}
