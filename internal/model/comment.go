package model

// 使用两层型
type FirstComment struct { //直接评论帖子的一级评论
	CommonModel
	BlogID       int64  `gorm:"not null" json:"blog_id"`           // 所属帖子 ID
	UserID       int64  `gorm:"not null" json:"user_id"`           // 评论者 ID
	Content      string `gorm:"type:text;not null" json:"content"` // 评论内容
	LikesCount   int    `gorm:"default:0" json:"likes_count"`      // 点赞数
	RepliesCount int    `gorm:"default:0" json:"replies_count"`    // 回复数
	// 关联二级评论
	SecondComments []SecondComment `gorm:"foreignKey:RootParentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (t *FirstComment) TableName() string {
	return "first_comment"
}

/*
对一级评论进行回复的二级评论
对二级评论进行回复的二级评论
*/
type SecondComment struct {
	CommonModel
	BlogID       int64  `gorm:"not null" json:"blog_id"`           // 所属帖子 ID
	UserID       int64  `gorm:"not null" json:"user_id"`           // 评论者 ID
	Content      string `gorm:"type:text;not null" json:"content"` // 评论内容
	ParentID     int64  `gorm:"default:0" json:"parent_id"`        // 父评论ID（为0是回复一级评论的评论，不为0是回复二级评论的评论）
	RootParentID int64  `gorm:"default:0" json:"root_parent_id"`   // 根评论ID（所属的级评论）
	RepliesCount int    `gorm:"default:0" json:"replies_count"`    // 回复数
}

func (t *SecondComment) TableName() string {
	return "second_comment"
}

type CommentLike struct {
	CommonModel
	CommentID int64 `gorm:"not null" json:"comment_id"` // 所属评论或回复ID
	UserID    int64 `gorm:"not null" json:"user_id"`    // 点赞用户ID
	IsLiked   bool  `gorm:"default:false"`
}
