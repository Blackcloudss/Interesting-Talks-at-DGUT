package repo

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"gorm.io/gorm"
)

// SearchRepo 搜索帖子仓库
type SearchRepo struct {
	DB *gorm.DB
}

// NewSearchRepo 创建搜索帖子仓库实例
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
			Joins("LEFT JOIN user_display ON blog.user_id = user_display.id").
			Where("blog.title LIKE ?", "%"+keyword+"%")
	case "tag":
		query = r.DB.Model(&model.Blog{}).
			Joins("LEFT JOIN user_display ON blog.user_id = user_display.id").
			Where("blog.blog_tag LIKE ? OR blog.sub_tag LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	case "nickname":
		query = r.DB.Model(&model.Blog{}).
			Joins("LEFT JOIN user_display ON blog.user_id = user_display.id").
			Where("user_display.nickname LIKE ?", "%"+keyword+"%")
	default: // 综合搜索
		query = r.DB.Model(&model.Blog{}).
			Joins("LEFT JOIN user_display ON blog.user_id = user_display.id").
			Where("blog.title LIKE ? OR blog.blog_tag LIKE ? OR blog.sub_tag LIKE ? OR user_display.nickname LIKE ?",
				"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 查询符合条件的帖子总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询帖子和用户信息
	if err := query.
		Select("blog.*, user_display.nickname, user_display.avatar, user_display.tag").
		Order("blog.created_at DESC").
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

// DeleteSearchHistory 删除指定的搜索历史记录
func (r *SearchRepo) DeleteSearchHistory(userID int64, historyID int64) error {
	// 删除指定的搜索历史记录
	return r.DB.Where("user_id = ? AND id = ?", userID, historyID).Delete(&model.SearchHistory{}).Error
}

// GetHotSearchRepo 获取热门搜索
func (r *SearchRepo) GetHotSearchRepo() (*types.GetHotSearchResp, error) {
	var hotSearchList []string
	if err := r.DB.Model(&model.Blog{}).
		Select("title").
		Order("created_at DESC, be_liked DESC").
		Limit(10).
		Scan(&hotSearchList).Error; err != nil {
		return nil, err
	}
	return &types.GetHotSearchResp{HotSearchList: hotSearchList}, nil
}
