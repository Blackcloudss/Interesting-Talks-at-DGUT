package repo

import (
	"errors"
	"fmt"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/pkg/redisx"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"gorm.io/gorm"
	"time"
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

// CheckBlogExists 检查帖子是否存在，如果存在则返回帖子对象
func (r *BlogRepo) CheckBlogExists(blogID int64) (*model.Blog, error) {
	var blog model.Blog
	// 查询帖子是否存在
	err := r.DB.Model(&model.Blog{}).Where("id = ?", blogID).First(&blog).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 如果帖子不存在，返回 nil
		}
		return nil, err // 如果发生其他错误，返回错误
	}
	return &blog, nil // 如果帖子存在，返回帖子对象
}

// GetBlogByID 根据ID获取帖子
func (r *BlogRepo) GetBlogByID(blogID int64) (types.BlogResp, error) {
	var blogResp types.BlogResp

	// 联合查询帖子和用户信息
	if err := r.DB.Model(&model.Blog{}).
		Select("blogs.id AS BlogID, blogs.updated_at AS UpdatedAt, blogs.be_liked AS BeLiked, blogs.be_collected AS BeCollected, blogs.comment_count AS CommentCount, blogs.tag AS BlogTag, blogs.sub_tag AS SubTag, blogs.view_permission AS ViewPermission, blogs.content AS Content, user_displays.user_id AS UserID, user_displays.nickname AS Nickname, user_displays.avatar AS Avatar, user_displays.tag AS Tag").
		Joins("LEFT JOIN user_displays ON blogs.user_id = user_displays.user_id").
		Where("blogs.id = ?", blogID).
		Scan(&blogResp).Error; err != nil {
		return types.BlogResp{}, err
	}

	return blogResp, nil
}

// GetBlogs 分页获取帖子列表
func (r *BlogRepo) GetBlogs(page, pageSize int) ([]types.BlogResp, int64, error) {
	var blogs []types.BlogResp
	var total int64

	// 查询帖子总数
	if err := r.DB.Model(&model.Blog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询帖子和用户信息
	if err := r.DB.Model(&model.Blog{}).
		Select("blogs.id AS BlogID, blogs.updated_at AS UpdatedAt, blogs.be_liked AS BeLiked, blogs.be_collected AS BeCollected, blogs.comment_count AS CommentCount, blogs.tag AS BlogTag, blogs.sub_tag AS SubTag, blogs.view_permission AS ViewPermission, blogs.content AS Content, user_displays.user_id AS UserID, user_displays.nickname AS Nickname, user_displays.avatar AS Avatar, user_displays.tag AS Tag").
		Joins("LEFT JOIN user_displays ON blogs.user_id = user_displays.user_id").
		Order("blogs.created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&blogs).Error; err != nil {
		return nil, 0, err
	}

	return blogs, total, nil
}

// GetBlogsByTag 根据标签获取帖子列表
func (r *BlogRepo) GetBlogsByTag(subTag string, page int, pageSize int) ([]types.BlogResp, int64, error) {
	var blogs []types.BlogResp
	var total int64

	// 查询符合条件的帖子总数
	if err := r.DB.Model(&model.Blog{}).
		Where("sub_tag = ?", subTag).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询帖子和用户信息
	if err := r.DB.Model(&model.Blog{}).
		Select("blogs.id AS BlogID, blogs.updated_at AS UpdatedAt, blogs.be_liked AS BeLiked, blogs.be_collected AS BeCollected, blogs.comment_count AS CommentCount, blogs.tag AS BlogTag, blogs.sub_tag AS SubTag, blogs.view_permission AS ViewPermission, blogs.content AS Content, user_displays.user_id AS UserID, user_displays.nickname AS Nickname, user_displays.avatar AS Avatar, user_displays.tag AS Tag").
		Joins("LEFT JOIN user_displays ON blogs.user_id = user_displays.user_id").
		Where("sub_tag = ?", subTag).
		Order("blogs.created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&blogs).Error; err != nil {
		return nil, 0, err
	}

	return blogs, total, nil
}

// GetBlogsByUserID 根据用户ID分页获取帖子列表
func (r *BlogRepo) GetBlogsByUserID(userID int64, page, pageSize int) ([]types.BlogResp, int64, error) {
	var blogs []types.BlogResp
	var total int64

	// 查询用户帖子总数
	if err := r.DB.Model(&model.Blog{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询帖子和用户信息
	if err := r.DB.Model(&model.Blog{}).
		Select("blogs.id AS BlogID, blogs.updated_at AS UpdatedAt, blogs.be_liked AS BeLiked, blogs.be_collected AS BeCollected, blogs.comment_count AS CommentCount, blogs.tag AS BlogTag, blogs.sub_tag AS SubTag, blogs.view_permission AS ViewPermission, blogs.content AS Content, user_displays.user_id AS UserID, user_displays.nickname AS Nickname, user_displays.avatar AS Avatar, user_displays.tag AS Tag").
		Joins("LEFT JOIN user_displays ON blogs.user_id = user_displays.user_id").
		Where("blogs.user_id = ?", userID).
		Order("blogs.created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&blogs).Error; err != nil {
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
	lockKey := fmt.Sprintf("blog:collect:%d", blogID)
	lockValue := fmt.Sprintf("%d-%d", userID, time.Now().UnixNano())
	expiration := 10 * time.Second

	// 尝试获取锁
	locked, err := redisx.Lock(global.Rdb, lockKey, lockValue, expiration)
	if err != nil {
		return err
	}
	if !locked {
		return errors.New("failed to acquire lock")
	}
	defer func() {
		// 释放锁
		if err := redisx.Unlock(global.Rdb, lockKey, lockValue); err != nil {
			zlog.Errorf("Failed to unlock: %v", err)
		}
	}()

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

// GetCollectedBlogs 获取用户收藏的帖子（分页）
func (r *BlogRepo) GetCollectedBlogs(userID int64, page, pageSize int) ([]types.BlogResp, int64, error) {
	var blogs []types.BlogResp
	var total int64

	// 查询用户收藏的帖子总数
	if err := r.DB.Model(&model.Blog{}).
		Joins("INNER JOIN collections ON collections.blog_id = blogs.id").
		Where("collections.user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询帖子和用户信息
	if err := r.DB.Model(&model.Blog{}).
		Select("blogs.id AS BlogID, blogs.updated_at AS UpdatedAt, blogs.be_liked AS BeLiked, blogs.be_collected AS BeCollected, blogs.comment_count AS CommentCount, blogs.tag AS BlogTag, blogs.sub_tag AS SubTag, blogs.view_permission AS ViewPermission, blogs.content AS Content, user_displays.user_id AS UserID, user_displays.nickname AS Nickname, user_displays.avatar AS Avatar, user_displays.tag AS Tag").
		Joins("INNER JOIN collections ON collections.blog_id = blogs.id").
		Joins("LEFT JOIN user_displays ON blogs.user_id = user_displays.user_id").
		Where("collections.user_id = ?", userID).
		Order("blogs.created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&blogs).Error; err != nil {
		return nil, 0, err
	}

	return blogs, total, nil
}

// IsBlogLiked 检查用户是否已经点赞某个帖子
func (r *BlogRepo) IsBlogLiked(userID, blogID int64) (bool, error) {
	var isLiked bool
	err := r.DB.Model(&model.Blog{}).Select("is_liked").Where("id = ? AND user_id = ?", blogID, userID).Scan(&isLiked).Error
	return isLiked, err
}

// LikeBlog 点赞帖子
func (r *BlogRepo) LikeBlog(userID, blogID int64) error {
	lockKey := fmt.Sprintf("blog:like:%d", blogID)
	lockValue := fmt.Sprintf("%d-%d", userID, time.Now().UnixNano())
	expiration := 10 * time.Second

	// 尝试获取锁
	locked, err := redisx.Lock(global.Rdb, lockKey, lockValue, expiration)
	if err != nil {
		return err
	}
	if !locked {
		return errors.New("failed to acquire lock")
	}
	defer func() {
		// 释放锁
		if err := redisx.Unlock(global.Rdb, lockKey, lockValue); err != nil {
			zlog.Errorf("Failed to unlock: %v", err)
		}
	}()

	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

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
