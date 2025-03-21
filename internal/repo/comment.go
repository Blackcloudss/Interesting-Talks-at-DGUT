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

// CreateFirstComment 创建一级评论并更新帖子的评论数
func (r *CommentRepo) CreateFirstComment(comment *model.FirstComment) error {
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

// CreateFirstComment 创建二级评论并更新帖子的评论数
func (r *CommentRepo) CreateSecondComment(comment *model.SecondComment) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// 对评论内容进行敏感词过滤替换
		comment.Content = global.Filter.Replace(comment.Content, '*')
		// 创建评论记录
		if err := tx.Create(comment).Error; err != nil {
			return err
		}
		// 更新父评论的回复数
		if err := tx.Model(&model.FirstComment{}).Where("id = ?", comment.ParentID).Update("replies_count", gorm.Expr("replies_count + 1")).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeleteComment 删除评论
func (r *CommentRepo) DeleteComment(commentID, blogID int64, isFirstComment bool) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if isFirstComment {
			// 删除一级评论记录
			if err := tx.Delete(&model.FirstComment{}, commentID).Error; err != nil {
				return err
			}
			// 更新帖子的评论数
			if err := tx.Model(&model.Blog{}).Where("id = ?", blogID).Update("comment_count", gorm.Expr("comment_count - 1")).Error; err != nil {
				return err
			}
		} else {
			// 删除二级评论记录
			if err := tx.Delete(&model.SecondComment{}, commentID).Error; err != nil {
				return err
			}
			// 获取二级评论的根评论ID（即父一级评论ID）
			var secondComment model.SecondComment
			if err := tx.Where("id = ?", commentID).First(&secondComment).Error; err != nil {
				return err
			}
			// 更新父一级评论的回复数
			if err := tx.Model(&model.FirstComment{}).Where("id = ?", secondComment.RootParentID).Update("replies_count", gorm.Expr("replies_count - 1")).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

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

	// 插入点赞记录
	if err := tx.Create(&model.CommentLike{
		CommentID: commentID,
		UserID:    userID,
		IsLiked:   true,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新评论的点赞数
	if err := tx.Model(&model.FirstComment{}).Where("id = ?", commentID).Update("likes_count", gorm.Expr("likes_count + 1")).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
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

	// 删除点赞记录
	if err := tx.Where("user_id = ? AND comment_id = ?", userID, commentID).Delete(&model.CommentLike{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新评论的点赞数
	if err := tx.Model(&model.FirstComment{}).Where("id = ?", commentID).Update("likes_count", gorm.Expr("likes_count - 1")).Error; err != nil {
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

// GetCommentByID 根据评论ID获取评论
func (r *CommentRepo) GetCommentByID(commentID int64) (*model.FirstComment, *model.SecondComment, error) {
	// 尝试从一级评论表中获取
	var firstComment model.FirstComment
	if err := r.DB.First(&firstComment, commentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果在一级评论表中未找到，尝试从二级评论表中获取
			var secondComment model.SecondComment
			if err := r.DB.First(&secondComment, commentID).Error; err != nil {
				zlog.Errorf("Failed to get comment by ID: %v", err)
				return nil, nil, err
			}
			// 找到二级评论，返回二级评论
			return nil, &secondComment, nil
		} else {
			zlog.Errorf("Failed to get comment by ID: %v", err)
			return nil, nil, err
		}
	}
	// 找到一级评论，返回一级评论
	return &firstComment, nil, nil
}

/*
// GetCommentList 获取评论列表
func (r *CommentRepo) GetCommentList(blogID int64) ([]model.Comment, error) {
	var comments []model.Comment
	if err := r.DB.Where("blog_id = ?", blogID).Order("created_at DESC").Find(&comments).Error; err != nil {
		zlog.Errorf("Failed to get comment list: %v", err)
		return nil, err
	}
	return comments, nil
}*/

/*
// GetRepliesList 获取评论的回复列表
func (r *CommentRepo) GetRepliesList(commentID int64) ([]model.Comment, error) {
	var replies []model.Comment
	err := r.DB.Where("parent_id = ?", commentID).Order("created_at DESC").Find(&replies).Error
	return replies, err
}*/
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
