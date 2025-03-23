package model

// 关注表 -- 获取好友列表需要
type Follow struct {
	CommonModel
	FollowerID int64 `gorm:"column:follower_id;type:bigint;comment:'关注者ID';index" json:"follower_id"`  // 关注者ID
	FollowedID int64 `gorm:"column:followed_id;type:bigint;comment:'被关注者ID';index" json:"followed_id"` // 被关注者ID
	IsFriend   bool  `gorm:"column:is_friend;type:tinyint(1);comment:'是否为好友'" json:"is_friend"`        // 是否为好友
}

func (Follow) TableName() string { return "follow" }
