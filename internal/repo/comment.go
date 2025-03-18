package repo

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"gorm.io/gorm"
)

type CommentRepo struct {
	DB *gorm.DB
}

func NewCommentRepo(db *gorm.DB) *CommentRepo {
	return &CommentRepo{
		DB: db,
	}
}

// CreateComment 创建评论并更新帖子的评论数
func (r *CommentRepo) CreateComment(comment *model.Comment) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// 对评论内容进行敏感词过滤替换
		comment.Content = global.Filter.Replace(comment.Content, '*')

		// 创建评论记录
		if err := tx.Create(comment).Error; err != nil {
			return err
		}

		// 更新帖子的评论数
		if err := tx.Model(&model.Blog{}).Where("id = ?", comment.BlogID).Update("comment_count", gorm.Expr("comment_count + 1")).Error; err != nil {
			return err
		}

		return nil
	})
}

// DeleteComment 删除评论并减少帖子的评论数
func (r *CommentRepo) DeleteComment(commentID, blogID int64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// 删除评论记录
		if err := tx.Delete(&model.Comment{}, commentID).Error; err != nil {
			return err
		}

		// 更新帖子的评论数
		if err := tx.Model(&model.Blog{}).Where("id = ?", blogID).Update("comment_count", gorm.Expr("comment_count - 1")).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetCommentList 获取评论列表
func (r *CommentRepo) GetCommentList(blogID int64) ([]model.Comment, error) {
	var comments []model.Comment
	if err := r.DB.Where("blog_id = ?", blogID).Order("created_at DESC").Find(&comments).Error; err != nil {
		zlog.Errorf("Failed to get comment list: %v", err)
		return nil, err
	}
	return comments, nil
}

// GetCommentByID 根据评论ID获取评论
func (r *CommentRepo) GetCommentByID(commentID int64) (*model.Comment, error) {
	var comment model.Comment
	if err := r.DB.First(&comment, commentID).Error; err != nil {
		zlog.Errorf("Failed to get comment by ID: %v", err)
		return nil, err
	}
	return &comment, nil
}
