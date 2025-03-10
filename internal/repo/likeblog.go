package repo

import (
	"errors"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"gorm.io/gorm"
)

// LikeBlog 点赞帖子
func LikeBlog(tx *gorm.DB, userID, blogID uint64) error {
	// 检查是否已经点赞
	var like model.Like
	if err := tx.Where("user_id = ? AND blog_id = ?", userID, blogID).First(&like).Error; err == nil {
		return errors.New("already liked this blog")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 创建点赞记录
	if err := tx.Create(&model.Like{
		UserID: userID,
		BlogID: blogID,
	}).Error; err != nil {
		return err
	}

	// 更新帖子的点赞数
	if err := tx.Model(&model.Blog{}).Where("id = ?", blogID).Update("like_count", gorm.Expr("like_count + 1")).Error; err != nil {
		return err
	}

	return nil
}

// UnlikeBlog 取消点赞
func UnlikeBlog(tx *gorm.DB, userID, blogID uint64) error {
	// 删除点赞记录
	if err := tx.Where("user_id = ? AND blog_id = ?", userID, blogID).Delete(&model.Like{}).Error; err != nil {
		return err
	}

	// 更新帖子的点赞数
	if err := tx.Model(&model.Blog{}).Where("id = ?", blogID).Update("like_count", gorm.Expr("like_count - 1")).Error; err != nil {
		return err
	}

	return nil
}

// IsLiked 检查是否已经点赞
func IsLiked(tx *gorm.DB, userID, blogID uint64) bool {
	var count int64
	tx.Model(&model.Like{}).Where("user_id = ? AND blog_id = ?", userID, blogID).Count(&count)
	return count > 0
}
