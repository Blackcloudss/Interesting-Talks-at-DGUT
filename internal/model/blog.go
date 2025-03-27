package model

type Blog struct {
	CommonModel
	LikeCount      int    `gorm:"column:like_count;default:0" json:"like_count"`
	CollectCount   int    `gorm:"column:collect_count;default:0" json:"collect_count"`
	CommentCount   int    `gorm:"column:comment_count;default:0" json:"comment_count"`
	BlogTag        string `gorm:"column:blog_tag;size:30;index" json:"blog_tag"`
	SubTag         string `gorm:"column:sub_tag;size:50;index" json:"sub_tag"`
	ViewPermission string `gorm:"column:view_permission;type:varchar(20);default:'所有人'" json:"view_permission"`
	UserID         int64  `gorm:"column:user_id" json:"user_id"`
	Title          string `gorm:"column:title;size:100;not null;index" json:"title"`
	Content        string `gorm:"column:content;type:text;not null" json:"content"`
}

func (t *Blog) TableName() string {
	return "blog"
}

type Like struct { // 点赞
	CommonModel
	UserID  int64 `gorm:"column:user_id" json:"user_id"`
	BlogID  int64 `gorm:"column:blog_id" json:"blog_id"`
	IsLiked bool  `gorm:"column:is_liked;default:false" json:"is_liked"` // 是否被点赞，默认为false
}

func (l *Like) TableName() string { return "like" }

type Collection struct {
	CommonModel
	UserID      int64 `gorm:"column:user_id" json:"user_id"`                         // 进行收藏操作的用户ID
	BlogID      int64 `gorm:"column:blog_id" json:"blog_id"`                         // 对应的BlogID
	IsCollected bool  `gorm:"column:is_collected;default:false" json:"is_collected"` // 是否被收藏，默认为false
}

func (c *Collection) TableName() string { return "collection" }
