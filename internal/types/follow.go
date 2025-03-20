package types

// FollowReq 关注请求结构体
type FollowReq struct {
	FollowedID int64 `json:"followedID"` // 被关注者ID
}

// FollowResp 关注响应结构体
type FollowResp struct {
}

// UnfollowReq 取消关注请求结构体
type UnfollowReq struct {
	FollowedID int64 `json:"followedID"` // 被关注者ID
}

// UnfollowResp 取消关注响应结构体
type UnfollowResp struct {
}

// FollowInfo 关注的用户或粉丝通用信息
type FollowInfo struct {
	UserID       int64  `json:"userID"`       // 用户ID
	Nickname     string `json:"nickname"`     // 昵称
	Avatar       string `json:"avatar"`       // 头像
	FollowedAt   string `json:"followedAt"`   // 已关注的时间
	IsFollowedBy bool   `json:"isFollowedBy"` // 当前用户是否关注了该用户
}

// GetFollowingsReq 获取我的关注列表请求结构体
type GetFollowingsReq struct {
}

// GetOtherFollowingsReq 获取其他用户关注列表请求结构体
type GetOtherFollowingsReq struct {
	OtherUserID int64 `json:"otherUserID"` // 其他用户ID
}

// 获取我的关注列表响应结构体
type GetFollowingsResp struct {
	Followings []FollowInfo `json:"followings"`
}

// GetFollowersReq 获取我的粉丝列表请求结构体
type GetFollowersReq struct {
}

// GetOtherFollowersReq 获取其他用户粉丝列表请求结构体
type GetOtherFollowersReq struct {
	OtherUserID int64 `json:"otherUserID"` // 其他用户ID
}

// 获取用户粉丝列表响应结构体
type GetFollowersResp struct {
	Followers []FollowInfo `json:"followers"`
}
