package model

import (
	"gorm.io/gorm"
	"time"
)

type Blog struct { // 帖子
	ID             uint64         `gorm:"primaryKey;autoIncrement:false" json:"id"`              // 帖子ID
	CreatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`           // 创建时间
	UpdatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`           // 更新时间
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`                               // 删除时间
	BeLiked        int            `gorm:"default:0" json:"be_liked"`                             // 点赞数，默认为0
	BeCollected    int            `gorm:"default:0" json:"be_collected"`                         // 收藏数，默认为0
	CommentCount   int            `gorm:"default:0" json:"comment_count"`                        // 评论数，默认为0（字段名更清晰）
	Tag            string         `gorm:"size:30" json:"tag"`                                    // 帖子分区（主标签）
	SubTag         string         `gorm:"size:50" json:"sub_tag"`                                // 子标签
	ViewPermission string         `gorm:"type:varchar(20);default:'所有人'" json:"view_permission"` // 访问权限，默认为“所有人”
	UserID         uint64         `json:"user_id"`                                               // 用户ID
	Content        string         `gorm:"type:text;not null" json:"content"`                     // 内容，不可为空
}

type Like struct { // 点赞
	UserID    uint64    `gorm:"primaryKey" json:"user_id"`
	BlogID    uint64    `gorm:"primaryKey" json:"blog_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

type Collection struct {
	UserID    uint64    `gorm:"primaryKey" json:"user_id"` // 进行收藏操作的用户ID
	BlogID    uint64    `gorm:"primaryKey" json:"blog_id"` // 对应的BlogID
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

type Comment struct {
	CommentID uint64         `gorm:"primaryKey;autoIncrement:false" json:"comment_id"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"` // 创建时间
	UpdatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"` // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`                     // 删除时间
	BlogID    uint64         `gorm:"not null" json:"blog_id"`                     // 所属帖子ID
	AuthorID  uint64         `gorm:"not null" json:"author_id"`                   // 评论作者ID（修正字段名）
	Content   string         `gorm:"type:text;not null" json:"content"`           // 评论内容
}
