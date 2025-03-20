package model

type Blog struct { // 帖子
	CommonModel           // 删除时间
	BeLiked        int    `gorm:"default:0" json:"be_liked"`                             // 点赞数，默认为0
	BeCollected    int    `gorm:"default:0" json:"be_collected"`                         // 收藏数，默认为0
	CommentCount   int    `gorm:"default:0" json:"comment_count"`                        // 评论数，默认为0（字段名更清晰）
	BlogTag        string `gorm:"size:30;index" json:"blog_tag"`                         // 帖子分区（主标签）
	SubTag         string `gorm:"size:50;index" json:"sub_tag"`                          // 子标签
	ViewPermission string `gorm:"type:varchar(20);default:'所有人'" json:"view_permission"` // 访问权限，默认为“所有人”
	UserID         int64  `json:"user_id"`                                               // 用户ID
	Title          string `gorm:"size:100;not null;index" json:"title"`                  //标题
	Content        string `gorm:"type:text;not null" json:"content"`                     // 内容，不可为空
}

type Like struct { // 点赞
	CommonModel
	UserID  int64 ` json:"user_id"`
	BlogID  int64 ` json:"blog_id"`
	IsLiked bool  `gorm:"default:false" json:"is_liked"` // 是否被点赞，默认为false
}

type Collection struct {
	CommonModel
	UserID      int64 ` json:"user_id"`                          // 进行收藏操作的用户ID
	BlogID      int64 ` json:"blog_id"`                          // 对应的BlogID
	IsCollected bool  `gorm:"default:false" json:"is_collected"` // 是否被收藏，默认为false
}
