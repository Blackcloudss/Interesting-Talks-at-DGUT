package model

// 历史搜索库
type SearchHistory struct {
	CommonModel
	UserID        int64  `gorm:"column:user_id;not null" json:"user_id"`
	SearchContent string `gorm:"column:search_content;type:varchar(255);not null" json:"search_content"`
}
