package model

type Notice struct {
	CommonModel        // 删除时间
	UserID      int64  `json:"user_id"`                           // 用户ID
	Content     string `gorm:"type:text;not null" json:"content"` // 内容，不可为空
}

func (t *Notice) TableName() string {
	return "notice"
}
