package repo

import (
	"fmt"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/pkg/redisx"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
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

/*
// GetRepliesList 获取评论的回复列表
func (r *CommentRepo) GetRepliesList(commentID int64) ([]model.Comment, error) {
	var replies []model.Comment
	err := r.DB.Where("parent_id = ?", commentID).Order("created_at DESC").Find(&replies).Error
	return replies, err
}*/

// GetCommentByID 根据评论ID获取评论
func (r *CommentRepo) GetCommentByID(commentID int64) (*model.Comment, error) {
	var comment model.Comment
	if err := r.DB.First(&comment, commentID).Error; err != nil {
		zlog.Errorf("Failed to get comment by ID: %v", err)
		return nil, err
	}
	return &comment, nil
}

/*
// GetCommentListWithReplies 获取评论列表及其回复
func (r *CommentRepo) GetCommentListWithReplies(blogID int64) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.DB.Where("blog_id = ? AND parent_id = 0", blogID).Order("created_at DESC").Find(&comments).Error
	if err != nil {
		return nil, err
	}

	for i := range comments {
		replies, err := r.GetRepliesList(comments[i].ID)
		if err != nil {
			return nil, err
		}
		comments[i].Replies = replies
	}

	return comments, nil
}*/

// LikeComment 点赞评论
func (r *CommentRepo) LikeComment(userID, commentID int64) error {
	lockKey := fmt.Sprintf("comment:like:%d", commentID)
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
	isLiked, err := r.IsCommentLiked(userID, commentID)
	if err != nil {
		tx.Rollback()
		return err
	}
	if isLiked {
		tx.Rollback()
		return errors.New("already liked this comment")
	}

	// 更新评论的点赞状态和点赞数
	if err := tx.Model(&model.Comment{}).Where("id = ?", commentID).Updates(map[string]interface{}{
		"is_liked":   true,
		"like_count": gorm.Expr("like_count + 1"),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// UnlikeComment 取消点赞评论
func (r *CommentRepo) UnlikeComment(userID, commentID int64) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 检查是否已经点赞
	isLiked, err := r.IsCommentLiked(userID, commentID)
	if err != nil {
		tx.Rollback()
		return err
	}
	if !isLiked {
		tx.Rollback()
		return errors.New("not liked this comment")
	}

	// 更新评论的点赞状态和点赞数
	if err := tx.Model(&model.Comment{}).Where("id = ?", commentID).Updates(map[string]interface{}{
		"is_liked":   false,
		"like_count": gorm.Expr("like_count - 1"),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// IsCommentLiked 检查用户是否已经点赞了评论
func (r *CommentRepo) IsCommentLiked(userID, commentID int64) (bool, error) {
	var like model.CommentLike
	result := r.DB.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&like)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}
	return like.IsLiked, nil
}
