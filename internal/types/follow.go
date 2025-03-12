package types

import "github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"

// FollowReq 关注请求结构体
type FollowReq struct {
	FollowerID int64 `json:"followerID"` // 关注者ID
	FollowedID int64 `json:"followedID"` // 被关注者ID
}

// FollowResp 关注响应结构体
type FollowResp struct {
	Success bool `json:"success"`
}

// UnfollowReq 取消关注请求结构体
type UnfollowReq struct {
	FollowerID int64 `json:"followerID"` // 关注者ID
	FollowedID int64 `json:"followedID"` // 被关注者ID
}

// UnfollowResp 取消关注响应结构体
type UnfollowResp struct {
	Success bool `json:"success"`
}

// GetFollowingsReq 获取关注列表请求结构体
type GetFollowingsReq struct {
	UserID int64 `json:"userID"` // 用户ID
}

// GetFollowingsResp 获取关注列表响应结构体
type GetFollowingsResp struct {
	List []model.User `json:"list"` // 关注列表
}

// GetFollowersReq 获取粉丝列表请求结构体
type GetFollowersReq struct {
	UserID int64 `json:"userID"` // 用户ID
}

// GetFollowersResp 获取粉丝列表响应结构体
type GetFollowersResp struct {
	List []model.User `json:"list"` // 粉丝列表
}
