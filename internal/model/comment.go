package model

type Comment struct {
	CommonModel
	BlogID       int64     `gorm:"not null" json:"blog_id"`           // 所属帖子 ID
	UserID       int64     `gorm:"not null" json:"user_id"`           // 评论作者 ID
	Content      string    `gorm:"type:text;not null" json:"content"` // 评论内容
	ParentID     int64     `gorm:"default:0" json:"parent_id"`        // 父评论 ID（=0为帖子的直接评论）
	LikesCount   int       `gorm:"default:0" json:"likes_count"`      // 点赞数
	RepliesCount int       `gorm:"default:0" json:"replies_count"`    // 回复数
	Replies      []Comment `gorm:"-" json:"replies"`                  // 回复列表（用于嵌套显示）
}

type CommentLike struct {
	CommonModel
	CommentID int64 `gorm:"not null" json:"comment_id"` // 所属评论ID
	UserID    int64 `gorm:"not null" json:"user_id"`    // 点赞用户ID
	IsLiked   bool  `gorm:"default:false"`
}
