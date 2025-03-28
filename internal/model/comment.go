package model

type Comment struct {
	CommonModel
	BlogID       int64  `gorm:"column:blog_id;not null" json:"blog_id"`
	UserID       int64  `gorm:"column:user_id;not null" json:"user_id"`
	Content      string `gorm:"column:content;type:text;not null" json:"content"`
	LikeCount    int    `gorm:"column:like_count;default:0" json:"like_count"`
	RepliesCount int    `gorm:"column:replies_count;default:0" json:"replies_count"`
	ParentID     int64  `gorm:"column:parent_id;default:null" json:"parent_id"`

	// 添加与帖子的关联关系
	Blog Blog `gorm:"foreignKey:BlogID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	// 自关联
	Parent *Comment `gorm:"foreignKey:ParentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (t *Comment) TableName() string {
	return "comment"
}

type CommentLike struct {
	CommonModel
	CommentID int64 `gorm:"column:comment_id;not null" json:"comment_id"` // 所属评论或回复ID
	UserID    int64 `gorm:"column:user_id;not null" json:"user_id"`       // 点赞用户ID
	IsLiked   bool  `gorm:"column:is_liked;default:false" json:"is_liked"`
}

func (t *CommentLike) TableName() string { return "comment_like" }
