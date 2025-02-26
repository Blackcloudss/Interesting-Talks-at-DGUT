package repo

import (
	"context"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"gorm.io/gorm"
)

// CreateComment 创建评论记录
func CreateComment(db *gorm.DB, comment *model.Comment) error {
	return db.Create(comment).Error
}

// DeleteComment 删除评论
func DeleteComment(db *gorm.DB, commentID uint64) error {
	return db.Delete(&model.Comment{}, commentID).Error
}

// IncrementCommentCount 增加帖子的评论数
func IncrementCommentCount(db *gorm.DB, blogID uint64) error {
	return db.Model(&model.Blog{}).Where("id = ?", blogID).Update("comment_count", gorm.Expr("comment_count + 1")).Error
}

// DecrementCommentCount 减少帖子的评论数
func DecrementCommentCount(db *gorm.DB, blogID uint64) error {
	return db.Model(&model.Blog{}).Where("id = ?", blogID).Update("comment_count", gorm.Expr("comment_count - 1")).Error
}

// GetCommentList 获取评论列表
func GetCommentList(ctx context.Context, blogID uint64) ([]model.Comment, error) {
	var comments []model.Comment
	if err := global.DB.WithContext(ctx).Where("blog_id = ?", blogID).Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// GetCommentByID 根据评论ID获取评论
func GetCommentByID(ctx context.Context, commentID uint64) (*model.Comment, error) {
	var comment model.Comment
	if err := global.DB.WithContext(ctx).First(&comment, commentID).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}
