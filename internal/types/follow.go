package types

import "github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"

// FollowReq 关注请求结构体
type FollowReq struct {
	FollowedID int64 `json:"followedID"` // 被关注者ID
}

// FollowResp 关注响应结构体
type FollowResp struct {
	Success bool `json:"success"`
}

// UnfollowReq 取消关注请求结构体
type UnfollowReq struct {
	FollowedID int64 `json:"followedID"` // 被关注者ID
}

// UnfollowResp 取消关注响应结构体
type UnfollowResp struct {
	Success bool `json:"success"`
}

// GetFollowingsReq 获取我的关注列表请求结构体
type GetFollowingsReq struct {
}

// GetFollowingsResp 获取关注列表响应结构体
type GetFollowingsResp struct {
	Followers []model.User `json:"followers"` // 关注列表
}

// GetFollowersReq 获取我的粉丝列表请求结构体
type GetFollowersReq struct {
}

// GetFollowersResp 获取粉丝列表响应结构体
type GetFollowersResp struct {
	Fans []model.User `json:"fans"` // 粉丝列表
}
