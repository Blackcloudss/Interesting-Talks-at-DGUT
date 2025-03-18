package model

type User struct {
	CommonModel
	Followers []User `gorm:"many2many:user_follows;joinForeignKey:ID;joinReferences:FollowedID" json:"followers"`
	Following []User `gorm:"many2many:user_follows;joinForeignKey:ID;joinReferences:FollowerID" json:"following"`
}

type Follow struct {
	CommonModel
	FollowerID int64 `gorm:"index" json:"follower_id"` // 关注者ID
	FollowedID int64 `gorm:"index" json:"followed_id"` // 被关注者ID
}
