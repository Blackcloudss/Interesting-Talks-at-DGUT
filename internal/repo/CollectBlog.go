package repo

import (
	"context"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"gorm.io/gorm"
)

// CollectBlog 收藏帖子
func CollectBlog(db *gorm.DB, userID, blogID uint64) error {
	return db.Create(&model.Collection{
		UserID: userID,
		BlogID: blogID,
	}).Error
}

// UncollectBlog 取消收藏
func UncollectBlog(db *gorm.DB, userID, blogID uint64) error {
	return db.Delete(&model.Collection{}, "user_id = ? AND blog_id = ?", userID, blogID).Error
}

// GetCollectedBlogs 获取用户收藏的帖子
func GetCollectedBlogs(ctx context.Context, userID uint64) ([]model.Blog, error) {
	var blogs []model.Blog
	err := global.DB.WithContext(ctx).Model(&model.Blog{}).
		Joins("inner join collections on collections.blog_id = blogs.id").
		Where("collections.user_id = ?", userID).
		Find(&blogs).Error
	return blogs, err
}

// IsCollected 检查是否已经收藏
func IsCollected(db *gorm.DB, userID, blogID uint64) bool {
	var count int64
	db.Model(&model.Collection{}).
		Where("user_id = ? AND blog_id = ?", userID, blogID).
		Count(&count)
	return count > 0
}
