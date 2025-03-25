package repo

import (
	"fmt"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/pkg/redisx"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
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

// CreateSecondComment 创建二级评论并更新帖子的评论数
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

// GetFirstCommentList 获取一级评论列表，并显示部分二级评论
func (r *CommentRepo) GetFirstCommentList(req types.GetCommentListReq) ([]types.CommentDetail, int64, error) {
	var (
		firstCommentDetails []types.CommentDetail
		totalCount          int64
	)

	// 查询一级评论总数
	if err := r.DB.Model(&model.FirstComment{}).Where("blog_id = ?", req.BlogID).Count(&totalCount).Error; err != nil {
		zlog.Errorf("查询一级评论总数失败：%v", err)
		return nil, 0, err
	}

	// 查询一级评论列表及用户信息
	if err := r.DB.Model(&model.FirstComment{}).
		Select("first_comments.*, user_display.nickname, user_display.avatar, user_display.tag").
		Joins("LEFT JOIN user_display ON first_comments.user_id = user_display.id").
		Where("first_comments.blog_id = ?", req.BlogID).
		Order(req.SortBy + " DESC").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Scan(&firstCommentDetails).Error; err != nil {
		zlog.Errorf("查询一级评论列表失败：%v", err)
		return nil, 0, err
	}

	// 提取一级评论的 ID 列表
	firstCommentIDs := make([]int64, len(firstCommentDetails))
	for i, fc := range firstCommentDetails {
		firstCommentIDs[i] = fc.FirstComment.ID
	}

	// 查询所有一级评论对应的二级评论（前3条）
	var secondCommentDetails []types.SecondCommentDetail
	if err := r.DB.Model(&model.SecondComment{}).
		Select("second_comments.*, user_display.nickname, user_display.avatar, user_display.tag").
		Joins("LEFT JOIN user_display ON second_comments.user_id = user_display.id").
		Where("second_comments.root_parent_id IN ?", firstCommentIDs).
		Order("second_comments.created_at ASC").
		Group("second_comments.root_parent_id").
		Limit(3).
		Scan(&secondCommentDetails).Error; err != nil {
		zlog.Warnf("查询部分二级评论失败：%v", err)
		secondCommentDetails = []types.SecondCommentDetail{}
	}

	// 构造二级评论映射
	secondCommentMap := make(map[int64][]types.SecondCommentDetail)
	for _, sc := range secondCommentDetails {
		secondCommentMap[sc.SecondComment.RootParentID] = append(secondCommentMap[sc.SecondComment.RootParentID], sc)
	}

	// 构造一级评论详情列表
	for i, fc := range firstCommentDetails {
		firstCommentDetails[i].SecondComments = secondCommentMap[fc.FirstComment.ID]
	}

	return firstCommentDetails, totalCount, nil
}

// GetSecondCommentList 获取更多二级评论
func (r *CommentRepo) GetSecondCommentList(req types.GetSecondCommentListReq) ([]types.SecondCommentDetail, int64, error) {
	var (
		secondCommentDetails []types.SecondCommentDetail
		totalCount           int64
	)

	// 查询二级评论总数
	if err := r.DB.Model(&model.SecondComment{}).
		Where("root_parent_id = ?", req.RootParentID).
		Count(&totalCount).Error; err != nil {
		zlog.Errorf("查询二级评论总数失败：%v", err)
		return nil, 0, err
	}

	// 查询二级评论列表及用户信息
	// 直接将结果映射到 types.SecondCommentDetail
	if err := r.DB.Model(&model.SecondComment{}).
		Select("second_comments.*, user_display.nickname, user_display.avatar, user_display.tag").
		Joins("LEFT JOIN user_display ON second_comments.user_id = user_display.id").
		Where("second_comments.root_parent_id = ?", req.RootParentID).
		Order("second_comments.created_at ASC").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Scan(&secondCommentDetails).Error; err != nil {
		zlog.Errorf("查询二级评论列表失败：%v", err)
		return nil, 0, err
	}

	return secondCommentDetails, totalCount, nil
}
