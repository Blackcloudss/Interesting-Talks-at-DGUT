package repo

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"gorm.io/gorm"
)

// BlogRepo 帖子仓库
type SearchRepo struct {
	DB *gorm.DB
}

// NewBlogRepo 创建帖子仓库实例
func NewSearchRepo(db *gorm.DB) *SearchRepo {
	return &SearchRepo{
		DB: db,
	}
}

// SearchBlogs 搜索帖子
func (r *SearchRepo) SearchBlogs(keyword string, searchType string, page, pageSize int) ([]types.BlogResp, int64, error) {
	var blogs []types.BlogResp
	var total int64

	var query *gorm.DB
	switch searchType {
	case "title":
		query = r.DB.Model(&model.Blog{}).
			Joins("LEFT JOIN user_displays ON blogs.user_id = user_displays.user_id").
			Where("blogs.title LIKE ?", "%"+keyword+"%")
	case "tag":
		query = r.DB.Model(&model.Blog{}).
			Joins("LEFT JOIN user_displays ON blogs.user_id = user_displays.user_id").
			Where("blogs.blog_tag LIKE ? OR blogs.sub_tag LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	case "nickname":
		query = r.DB.Model(&model.Blog{}).
			Joins("LEFT JOIN user_displays ON blogs.user_id = user_displays.user_id").
			Where("user_displays.nickname LIKE ?", "%"+keyword+"%")
	default: // 综合搜索
		query = r.DB.Model(&model.Blog{}).
			Joins("LEFT JOIN user_displays ON blogs.user_id = user_displays.user_id").
			Where("blogs.title LIKE ? OR blogs.blog_tag LIKE ? OR blogs.sub_tag LIKE ? OR user_displays.nickname LIKE ?",
				"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 查询符合条件的帖子总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询帖子和用户信息
	if err := query.
		Select("blogs.*, user_displays.nickname, user_displays.avatar, user_displays.tag").
		Order("blogs.created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&blogs).Error; err != nil {
		return nil, 0, err
	}

	return blogs, total, nil
}

// RecordSearchHistory 记录搜索历史
func (r *SearchRepo) RecordSearchHistory(userID int64, searchContent string) error {
	searchHistory := model.SearchHistory{
		UserID:        userID,
		SearchContent: searchContent,
	}

	return r.DB.Create(&searchHistory).Error
}

// GetSearchHistory 获取用户的搜索历史
func (r *SearchRepo) GetSearchHistory(userID int64, page, pageSize int) ([]model.SearchHistory, int64, error) {
	var history []model.SearchHistory
	var total int64

	// 查询用户搜索历史总数
	if err := r.DB.Model(&model.SearchHistory{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询搜索历史
	if err := r.DB.Model(&model.SearchHistory{}).
		Where("user_id = ?", userID).
		Order("search_time DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&history).Error; err != nil {
		return nil, 0, err
	}

	return history, total, nil
}
