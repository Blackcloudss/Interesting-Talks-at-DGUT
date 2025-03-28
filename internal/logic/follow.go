package logic

import (
	"context"
	"fmt"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/repo"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils"
	"time"
)

// 定义内部逻辑错误
var (
	codeFollowFailed        = response.MsgCode{Code: 40035, Msg: "关注或取消关注失败"}
	codeGetFollowingsFailed = response.MsgCode{Code: 40039, Msg: "获取关注列表失败"}
	codeGetFollowersFailed  = response.MsgCode{Code: 40040, Msg: "获取粉丝列表失败"}
)

type FollowLogic struct{}

func NewFollowLogic() *FollowLogic {
	return &FollowLogic{}
}

// FollowLogic.go
func (l *FollowLogic) Follow(ctx context.Context, req types.FollowReq, UserID int64) (resp *types.FollowResp, err error) {
	defer utils.RecordTime(time.Now())()

	// 用 redis 加锁
	lockKey := fmt.Sprintf("follow:lock:follower:%d:followed:%d", UserID, req.FollowedID)
	locked, err := global.Rdb.SetNX(ctx, lockKey, 1, 1*time.Second).Result()
	if err != nil {
		zlog.CtxErrorf(ctx, "Redis 上锁失败: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}
	if !locked {
		zlog.CtxInfof(ctx, "关注/取消关注操作正被 user: %d, followed: %d 使用，请稍等 1 s", UserID, req.FollowedID)
		return nil, response.ErrResp(err, response.USER_OPERATION_LOCKED)
	}
	defer global.Rdb.Del(ctx, lockKey)

	resp, err = repo.NewFollowRepo(global.DB).Follow(UserID, req.FollowedID)
	if err != nil {
		zlog.CtxErrorf(ctx, "ToggleFollow failed: %v", err)
		return nil, response.ErrResp(err, codeFollowFailed)
	}
	return resp, nil
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
			UserID:      user.UserID,
			Nickname:    user.Nickname,
			Avatar:      user.Avatar,
			FollowedAt:  user.FollowedAt,
			IsFollowing: repo.NewFollowRepo(global.DB).IsFollowing(user.UserID, UserID),
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
			UserID:      user.UserID,
			Nickname:    user.Nickname,
			Avatar:      user.Avatar,
			FollowedAt:  user.FollowedAt,
			IsFollowing: repo.NewFollowRepo(global.DB).IsFollowing(user.UserID, UserID),
		}
	}

	zlog.CtxInfof(ctx, "Followers retrieved successfully (userID: %d)", UserID)
	return resp, nil
}
