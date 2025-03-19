package model

type Comment struct {
	CommonModel
	BlogID       int64  `gorm:"not null" json:"blog_id"`           // 所属帖子ID
	UserID       int64  `gorm:"not null" json:"user_id"`           // 评论作者ID
	Content      string `gorm:"type:text;not null" json:"content"` // 评论内容
	ParentID     int64  `gorm:"default:0" json:"parent_id"`        // 父评论ID（顶级评论，即帖子的直接评论为0）
	LikesCount   int    `gorm:"default:0" json:"likes_count"`      // 点赞数
	RepliesCount int    `gorm:"default:0" json:"replies_count"`    // 回复数
}
type CommentLike struct {
	CommonModel
	CommentID int64 `gorm:"not null" json:"comment_id"` // 所属评论ID
	UserID    int64 `gorm:"not null" json:"user_id"`    // 点赞用户ID
}
