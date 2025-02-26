package repo

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"gorm.io/gorm"
)

// CreateBlog 创建帖子
func CreateBlog(db *gorm.DB, blog *model.Blog) error {
	return db.Create(blog).Error
}

// UpdateBlog 更新帖子
func UpdateBlog(blog *model.Blog) error {
	return global.DB.Save(blog).Error
}

// DeleteBlog 删除帖子
func DeleteBlog(blogID uint64) error {
	return global.DB.Delete(&model.Blog{}, blogID).Error
}

// GetBlogByID 根据ID获取帖子
func GetBlogByID(blogID uint64) (*model.Blog, error) {
	var blog model.Blog
	if err := global.DB.First(&blog, blogID).Error; err != nil {
		return nil, err
	}
	return &blog, nil
}

// GetBlogs 分页获取帖子列表
func GetBlogs(page, pageSize int) ([]model.Blog, int64, error) {
	var blogs []model.Blog
	var total int64

	// 查询帖子总数
	if err := global.DB.Model(&model.Blog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询帖子列表
	if err := global.DB.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&blogs).Error; err != nil {
		return nil, 0, err
	}

	return blogs, total, nil
}

// GetBlogsAfterID 获取比指定 ID 更新的帖子
func GetBlogsAfterID(latestID uint64, pageSize int) ([]model.Blog, error) {
	var blogs []model.Blog

	// 查询比指定 ID 更新的帖子
	if err := global.DB.Where("id > ?", latestID).
		Order("created_at DESC").
		Limit(pageSize).
		Find(&blogs).Error; err != nil {
		return nil, err
	}

	return blogs, nil
}

// GetBlogsByTag 根据标签分页获取帖子列表
func GetBlogsByTag(tag string, page, pageSize int) ([]model.Blog, int64, error) {
	var blogs []model.Blog
	var total int64

	if err := global.DB.Model(&model.Blog{}).Where("tag = ?", tag).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := global.DB.Where("tag = ?", tag).Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&blogs).Error; err != nil {
		return nil, 0, err
	}

	return blogs, total, nil
}

// GetBlogsByUserID 根据用户ID分页获取帖子列表
func GetBlogsByUserID(userID uint64, page, pageSize int) ([]model.Blog, int64, error) {
	var blogs []model.Blog
	var total int64

	if err := global.DB.Model(&model.Blog{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := global.DB.Where("user_id = ?", userID).Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&blogs).Error; err != nil {
		return nil, 0, err
	}

	return blogs, total, nil
}
