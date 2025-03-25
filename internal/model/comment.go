package model

// 使用两层型
type FirstComment struct {
	CommonModel
	BlogID         int64           `gorm:"column:blog_id;not null" json:"blog_id"`
	Blog           Blog            `gorm:"foreignKey:BlogID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID         int64           `gorm:"column:user_id;not null" json:"user_id"`
	Content        string          `gorm:"column:content;type:text;not null" json:"content"`
	LikesCount     int             `gorm:"column:likes_count;default:0" json:"likes_count"`
	RepliesCount   int             `gorm:"column:replies_count;default:0" json:"replies_count"`
	SecondComments []SecondComment `gorm:"foreignKey:RootParentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKeyName:fk_first_comment_second_comment"`
}

func (t *FirstComment) TableName() string {
	return "first_comment"
}

type SecondComment struct {
	CommonModel
	BlogID       int64  `gorm:"column:blog_id;not null" json:"blog_id"`
	Blog         Blog   `gorm:"foreignKey:BlogID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID       int64  `gorm:"column:user_id;not null" json:"user_id"`
	Content      string `gorm:"column:content;type:text;not null" json:"content"`
	ParentID     int64  `gorm:"column:parent_id;default:0" json:"parent_id"`
	RootParentID int64  `gorm:"column:root_parent_id;default:0" json:"root_parent_id"`
	RepliesCount int    `gorm:"column:replies_count;default:0" json:"replies_count"`
}

func (t *SecondComment) TableName() string {
	return "second_comment"
}

type CommentLike struct {
	CommonModel
	CommentID int64 `gorm:"column:comment_id;not null" json:"comment_id"` // 所属评论或回复ID
	UserID    int64 `gorm:"column:user_id;not null" json:"user_id"`       // 点赞用户ID
	IsLiked   bool  `gorm:"column:is_liked;default:false"`
}

func (t *CommentLike) TableName() string { return "comment_like" }
