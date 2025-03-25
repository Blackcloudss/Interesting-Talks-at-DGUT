package model

// 使用两层型
type FirstComment struct {
	CommonModel
	BlogID       int64  `gorm:"column:blog_id;not null;constraint:fk_blog_first_comments:blog_id REFERENCES blog(id) ON DELETE CASCADE ON UPDATE CASCADE" json:"blog_id"` // 所属帖子 ID
	UserID       int64  `gorm:"column:user_id;not null" json:"user_id"`                                                                                                   // 评论者 ID
	Content      string `gorm:"column:content;type:text;not null" json:"content"`                                                                                         // 评论内容
	LikesCount   int    `gorm:"column:likes_count;default:0" json:"likes_count"`                                                                                          // 点赞数
	RepliesCount int    `gorm:"column:replies_count;default:0" json:"replies_count"`                                                                                      // 回复数
	// 关联二级评论
	SecondComments []SecondComment `gorm:"foreignKey:root_parent_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
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
	BlogID       int64  `gorm:"column:blog_id;not null;constraint:fk_blog_second_comments:blog_id REFERENCES blog(id) ON DELETE CASCADE ON UPDATE CASCADE" json:"blog_id"` // 所属帖子 ID
	UserID       int64  `gorm:"column:user_id;not null" json:"user_id"`                                                                                                    // 评论者 ID
	Content      string `gorm:"column:content;type:text;not null" json:"content"`                                                                                          // 评论内容
	ParentID     int64  `gorm:"column:parent_id;default:0" json:"parent_id"`                                                                                               // 父评论ID
	RootParentID int64  `gorm:"column:root_parent_id;default:0" json:"root_parent_id"`                                                                                     // 根评论ID
	RepliesCount int    `gorm:"column:replies_count;default:0" json:"replies_count"`                                                                                       // 回复数
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
