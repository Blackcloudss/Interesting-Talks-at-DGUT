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
		zlog.CtxErrorf(ctx, "Follow failed: %v", err)
		return nil, response.ErrResp(err, codeFollowFailed)
	}
	return resp, nil
}

// GetFollowings 获取用户关注的用户列表
func (l *FollowLogic) GetFollowings(ctx context.Context, currentUserID, targetUserID int64) (*types.GetFollowingsResp, error) {
	followings, err := repo.NewFollowRepo(global.DB).GetFollowings(targetUserID)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get followings for user (userID: %d): %v", targetUserID, err)
		return nil, response.ErrResp(err, codeGetFollowingsFailed)
	}

	// 构造响应
	resp := &types.GetFollowingsResp{
		Followings: make([]types.FollowInfo, len(followings)),
	}
	for i, user := range followings {
		// 如果是查询自己的关注列表，则 IsFollowing = true（因为是自己的关注列表）
		// 如果是查询别人的关注列表，则需要检查当前用户是否关注了这些用户
		isFollowing := targetUserID == currentUserID || repo.NewFollowRepo(global.DB).IsFollowing(currentUserID, user.UserID)
		resp.Followings[i] = types.FollowInfo{
			UserID:      user.UserID,
			Nickname:    user.Nickname,
			Avatar:      user.Avatar,
			FollowedAt:  user.FollowedAt,
			IsFollowing: isFollowing,
		}
	}

	zlog.CtxInfof(ctx, "Followings retrieved successfully (userID: %d)", targetUserID)
	return resp, nil
}

// GetFollowers 获取用户的粉丝列表
func (l *FollowLogic) GetFollowers(ctx context.Context, currentUserID, targetUserID int64) (*types.GetFollowersResp, error) {
	followers, err := repo.NewFollowRepo(global.DB).GetFollowers(targetUserID)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get followers for user (userID: %d): %v", targetUserID, err)
		return nil, response.ErrResp(err, codeGetFollowersFailed)
	}

	// 构造响应
	resp := &types.GetFollowersResp{
		Followers: make([]types.FollowInfo, len(followers)),
	}
	for i, user := range followers {
		// 检查当前用户是否关注了这些粉丝
		isFollowing := repo.NewFollowRepo(global.DB).IsFollowing(currentUserID, user.UserID)
		resp.Followers[i] = types.FollowInfo{
			UserID:      user.UserID,
			Nickname:    user.Nickname,
			Avatar:      user.Avatar,
			FollowedAt:  user.FollowedAt,
			IsFollowing: isFollowing,
		}
	}

	zlog.CtxInfof(ctx, "Followers retrieved successfully (userID: %d)", targetUserID)
	return resp, nil
}
