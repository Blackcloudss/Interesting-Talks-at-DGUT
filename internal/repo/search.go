package repo

import (
	"fmt"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"gorm.io/gorm"
)

// SearchRepo 搜索仓库
type SearchRepo struct {
	DB *gorm.DB
}

// NewSearchRepo 创建搜索仓库实例
func NewSearchRepo(db *gorm.DB) *SearchRepo {
	return &SearchRepo{
		DB: db,
	}
}

const HOT_BLOGS_LIMIT = 10

// SearchBlogs 搜索帖子
func (r *SearchRepo) SearchBlogs(keyword string, searchType string, page, pageSize int) ([]types.BlogResp, int64, error) {
	var blogs []types.BlogResp
	var total int64

	query := r.DB.Model(&model.Blog{}).
		Select(BLOG_SELECT_FIELDS).
		Joins("LEFT JOIN user_display ON blog.user_id = user_display.id").
		Where(fmt.Sprintf("blog.%s IS NULL", DELETED_AT))

	// 根据搜索类型添加不同的查询条件
	switch searchType {
	case "title":
		query = query.Where("blog.title LIKE ?", "%"+keyword+"%")
	case "tag":
		query = query.Where("blog.blog_tag LIKE ? OR blog.sub_tag LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	case "nickname":
		query = query.Where("user_display.nickname LIKE ?", "%"+keyword+"%")
	default: // 综合搜索
		query = query.Where("blog.title LIKE ? OR blog.blog_tag LIKE ? OR blog.sub_tag LIKE ? OR user_display.nickname LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		zlog.Errorf("统计搜索结果总数失败: %v", err)
		return nil, 0, err
	}

	// 分页查询
	if err := query.
		Order(fmt.Sprintf("blog.%s DESC", CREATED_AT)).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&blogs).
		Error; err != nil {
		zlog.Errorf("搜索帖子失败: %v", err)
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
		Order("created_at DESC"). // 修改为按 created_at 排序
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

// GetHotSearchRepo 获取热搜榜
func (r *SearchRepo) GetHotSearchRepo() (*types.GetHotSearchResp, error) {
	var hotSearchList []string
	if err := r.DB.Model(&model.Blog{}).
		Select("title").
		Order("created_at DESC, like_count DESC").
		Limit(HOT_BLOGS_LIMIT).
		Scan(&hotSearchList).Error; err != nil {
		return nil, err
	}
	return &types.GetHotSearchResp{HotSearchList: hotSearchList}, nil
}
