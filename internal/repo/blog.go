package repo

import (
	"errors"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"gorm.io/gorm"
)

// BlogRepo 帖子仓库
type BlogRepo struct {
	DB *gorm.DB
}

// NewBlogRepo 创建帖子仓库实例
func NewBlogRepo(db *gorm.DB) *BlogRepo {
	return &BlogRepo{
		DB: db,
	}
}

// CreateBlog 创建帖子
func (r *BlogRepo) CreateBlog(blog *model.Blog) error {
	return r.DB.Create(blog).Error
}

// UpdateBlog 更新帖子
func (r *BlogRepo) UpdateBlog(blog *model.Blog) error {
	return r.DB.Save(blog).Error
}

// DeleteBlog 删除帖子
func (r *BlogRepo) DeleteBlog(blogID int64) error {
	return r.DB.Delete(&model.Blog{}, blogID).Error
}

// GetBlogByID 根据ID获取帖子
func (r *BlogRepo) GetBlogByID(blogID int64) (*model.Blog, error) {
	var blog model.Blog
	if err := r.DB.First(&blog, blogID).Error; err != nil {
		return nil, err
	}
	return &blog, nil
}

// GetBlogs 分页获取帖子列表
func (r *BlogRepo) GetBlogs(page, pageSize int) ([]model.Blog, int64, error) {
	var blogs []model.Blog
	var total int64

	if err := r.DB.Model(&model.Blog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.DB.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&blogs).Error; err != nil {
		return nil, 0, err
	}

	return blogs, total, nil
}

// GetBlogsByTag 根据标签获取帖子列表
func (r *BlogRepo) GetBlogsByTag(subTag string) ([]model.Blog, error) {
	var blogs []model.Blog

	// 查询符合条件的帖子
	if err := r.DB.Model(&model.Blog{}).Where("sub_tag = ?", subTag).Order("created_at DESC").Find(&blogs).Error; err != nil {
		return nil, err
	}

	return blogs, nil
}

// GetBlogsByUserID 根据用户ID分页获取帖子列表
func (r *BlogRepo) GetBlogsByUserID(userID int64, page, pageSize int) ([]model.Blog, int64, error) {
	var blogs []model.Blog
	var total int64

	if err := r.DB.Model(&model.Blog{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.DB.Where("user_id = ?", userID).Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&blogs).Error; err != nil {
		return nil, 0, err
	}

	return blogs, total, nil
}

// IsCollected 检查用户是否已经收藏某个帖子
func (r *BlogRepo) IsCollected(userID, blogID int64) (bool, error) {
	var count int64
	err := r.DB.Model(&model.Collection{}).
		Where("user_id = ? AND blog_id = ?", userID, blogID).
		Count(&count).Error
	return count > 0, err
}

// CollectBlog 用户收藏帖子
func (r *BlogRepo) CollectBlog(userID, blogID int64) error {
	// 检查是否已经收藏
	isCollected, err := r.IsCollected(userID, blogID)
	if err != nil {
		return err
	}
	if isCollected {
		return errors.New("already collected this blog")
	}

	// 创建新的收藏记录
	collection := model.Collection{
		UserID:      userID,
		BlogID:      blogID,
		IsCollected: true,
	}

	// 插入收藏记录
	if err := r.DB.Create(&collection).Error; err != nil {
		return err
	}

	// 更新帖子的收藏数
	return r.DB.Model(&model.Blog{}).Where("id = ?", blogID).Update("be_collected", gorm.Expr("be_collected + 1")).Error
}

// UncollectBlog 用户取消收藏帖子
func (r *BlogRepo) UncollectBlog(userID, blogID int64) error {
	// 检查是否已经收藏
	isCollected, err := r.IsCollected(userID, blogID)
	if err != nil {
		return err
	}
	if !isCollected {
		return errors.New("not collected this blog")
	}

	// 删除收藏记录
	if err := r.DB.Delete(&model.Collection{}, "user_id = ? AND blog_id = ?", userID, blogID).Error; err != nil {
		return err
	}

	// 更新帖子的收藏数
	return r.DB.Model(&model.Blog{}).Where("id = ?", blogID).Update("be_collected", gorm.Expr("be_collected - 1")).Error
}

// GetCollectedBlogs 获取用户收藏的帖子
func (r *BlogRepo) GetCollectedBlogs(userID int64) ([]model.Blog, error) {
	var blogs []model.Blog
	err := r.DB.Model(&model.Blog{}).
		Joins("inner join collections on collections.blog_id = blogs.id").
		Where("collections.user_id = ?", userID).
		Find(&blogs).Error
	return blogs, err
}

// IsBlogLiked 检查用户是否已经点赞某个帖子
func (r *BlogRepo) IsBlogLiked(userID, blogID int64) (bool, error) {
	var isLiked bool
	err := r.DB.Model(&model.Blog{}).Select("is_liked").Where("id = ? AND user_id = ?", blogID, userID).Scan(&isLiked).Error
	return isLiked, err
}

// LikeBlog 点赞帖子
func (r *BlogRepo) LikeBlog(userID, blogID int64) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 检查是否已经点赞
	isLiked, err := r.IsBlogLiked(userID, blogID)
	if err != nil {
		tx.Rollback()
		return err
	}
	if isLiked {
		tx.Rollback()
		return errors.New("already liked this blog")
	}

	// 更新帖子的点赞状态和点赞数
	if err := tx.Model(&model.Blog{}).Where("id = ?", blogID).Updates(map[string]interface{}{
		"is_liked":   true,
		"like_count": gorm.Expr("like_count + 1"),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// UnlikeBlog 取消点赞
func (r *BlogRepo) UnlikeBlog(userID, blogID int64) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 检查是否已经点赞
	isLiked, err := r.IsBlogLiked(userID, blogID)
	if err != nil {
		tx.Rollback()
		return err
	}
	if !isLiked {
		tx.Rollback()
		return errors.New("not liked this blog")
	}

	// 更新帖子的点赞状态和点赞数
	if err := tx.Model(&model.Blog{}).Where("id = ?", blogID).Updates(map[string]interface{}{
		"is_liked":   false,
		"like_count": gorm.Expr("like_count - 1"),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
