package types

// @Title        friend.go
// @Description
// @Create       XdpCs 2025-03-20 上午1:17
// @Update       XdpCs 2025-03-20 上午1:17

// 获取好友列表 入参
type GetFriendListReq struct {
}

// 获取好友列表 出参
type GetFriendListResp struct {
	Friends []FriendInfo
}

// 好友信息
type FriendInfo struct {
	ID       int64  `json:"id"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
	Tag      string `json:"tag"`
}
