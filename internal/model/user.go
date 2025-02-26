package model

type User struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"user_id"`
	Username  string `gorm:"unique;not null" json:"username"`          // 用户名唯一且不能为空
	Password  string `gorm:"not null" json:"password"`                 // 密码不能为空
	Role      string `gorm:"default:student" json:"role"`              // 默认角色为学生
	Nickname  string `gorm:"default:游客" json:"nickname"`               // 默认昵称为“游客”
	Avatar    string `gorm:"default:default_avatar.png" json:"avatar"` // 默认头像
	Followers []User `gorm:"many2many:user_follows;joinForeignKey:ID;joinReferences:FollowedID" json:"followers"`
	Following []User `gorm:"many2many:user_follows;joinForeignKey:ID;joinReferences:FollowerID" json:"following"`
}

type Follow struct {
	FollowerID uint64 `gorm:"primaryKey;autoIncrement:false" json:"follower_id"` // 关注者ID
	FollowedID uint64 `gorm:"primaryKey;autoIncrement:false" json:"followed_id"` // 被关注者ID
}
