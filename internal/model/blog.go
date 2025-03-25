package model

type Blog struct { // 帖子
	CommonModel           // 删除时间
	BeLiked        int    `gorm:"column:be_liked;default:0" json:"be_liked"  `                                  // 点赞数，默认为0
	BeCollected    int    `gorm:"column:be_collected;default:0" json:"be_collected"`                            // 收藏数，默认为0
	CommentCount   int    `gorm:"column:comment_count;default:0" json:"comment_count"`                          // 评论数，默认为0（字段名更清晰）
	BlogTag        string `gorm:"column:blog_tag;size:30;index" json:"blog_tag"`                                // 帖子分区（主标签）
	SubTag         string `gorm:"column:sub_tag;size:50;index" json:"sub_tag"`                                  // 子标签
	ViewPermission string `gorm:"column:view_permission;type:varchar(20);default:'所有人'" json:"view_permission"` // 访问权限，默认为“所有人”
	UserID         int64  `gorm:"column:user_id" json:"user_id"`                                                // 用户ID
	Title          string `gorm:"column:title;size:100;not null;index" json:"title"`                            // 标题
	Content        string `gorm:"column:content;type:text;not null" json:"content"`                             // 内容，不可为空
	// 关联一级评论
	FirstComments []FirstComment `gorm:"foreignKey:blog_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	// 关联二级评论
	SecondComments []SecondComment `gorm:"foreignKey:blog_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (b *Blog) TableName() string { return "blog" }

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
