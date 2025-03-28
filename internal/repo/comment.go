package repo

import (
	"fmt"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type CommentRepo struct {
	DB *gorm.DB
}

const (
	COMMENT_ID    = "comment_id"
	PARENT_ID     = "parent_id"
	REPLIES_COUNT = "replies_count"
)

func NewCommentRepo(db *gorm.DB) *CommentRepo {
	return &CommentRepo{
		DB: db,
	}
}

// CreateFirstComment 创建一级评论并更新帖子的评论数
func (r *CommentRepo) CreateFirstComment(comment *model.FirstComment) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		zlog.Errorf("创建一级评论事务开始失败: %v", tx.Error)
		return tx.Error
	}

	// 对评论内容进行敏感词过滤替换
	comment.Content = global.Filter.Replace(comment.Content, '*')

	// 创建评论记录
	if err := tx.Create(comment).Error; err != nil {
		tx.Rollback()
		zlog.Errorf("创建一级评论记录失败: %v", err)
		return err
	}

	// 更新帖子的评论数
	if err := tx.Model(&model.Blog{}).
		Where(fmt.Sprintf("%s = ?", BLOG_ID), comment.BlogID).
		Update(COMMENT_COUNT, gorm.Expr("comment_count + 1")).
		Error; err != nil {
		tx.Rollback()
		zlog.Errorf("更新帖子评论数失败: %v", err)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		zlog.Errorf("创建一级评论事务提交失败: %v", err)
		return err
	}

	return nil
}

// CreateSecondComment 创建二级评论并更新父一级评论的回复数
func (r *CommentRepo) CreateSecondComment(comment *model.SecondComment) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		zlog.Errorf("创建二级评论事务开始失败: %v", tx.Error)
		return tx.Error
	}

	// 对评论内容进行敏感词过滤替换
	comment.Content = global.Filter.Replace(comment.Content, '*')

	// 创建评论记录
	if err := tx.Create(comment).Error; err != nil {
		tx.Rollback()
		zlog.Errorf("创建二级评论记录失败: %v", err)
		return err
	}

	// 更新父一级评论的回复数
	if err := tx.Model(&model.FirstComment{}).
		Where(fmt.Sprintf("%s = ?", COMMENT_ID), comment.ParentID).
		Update(REPLIES_COUNT, gorm.Expr("replies_count + 1")).
		Error; err != nil {
		tx.Rollback()
		zlog.Errorf("更新一级评论回复数失败: %v", err)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		zlog.Errorf("创建二级评论事务提交失败: %v", err)
		return err
	}

	return nil
}

// DeleteComment 删除评论并更新相关计数
func (r *CommentRepo) DeleteComment(commentID, blogID int64, isFirstComment bool) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		zlog.Errorf("删除评论事务开始失败: %v", tx.Error)
		return tx.Error
	}

	if isFirstComment {
		// 删除一级评论记录
		if err := tx.Delete(&model.FirstComment{}, commentID).Error; err != nil {
			tx.Rollback()
			zlog.Errorf("删除一级评论失败: %v", err)
			return err
		}

		// 更新帖子的评论数
		if err := tx.Model(&model.Blog{}).
			Where(fmt.Sprintf("%s = ?", BLOG_ID), blogID).
			Update(COMMENT_COUNT, gorm.Expr("comment_count - 1")).
			Error; err != nil {
			tx.Rollback()
			zlog.Errorf("更新帖子评论数失败: %v", err)
			return err
		}
	} else {
		// 获取二级评论的父一级评论ID
		var secondComment model.SecondComment
		if err := tx.Where(fmt.Sprintf("%s = ?", COMMENT_ID), commentID).
			First(&secondComment).
			Error; err != nil {
			tx.Rollback()
			zlog.Errorf("查询二级评论失败: %v", err)
			return err
		}

		// 删除二级评论记录
		if err := tx.Delete(&model.SecondComment{}, commentID).Error; err != nil {
			tx.Rollback()
			zlog.Errorf("删除二级评论失败: %v", err)
			return err
		}

		// 更新父一级评论的回复数
		if err := tx.Model(&model.FirstComment{}).
			Where(fmt.Sprintf("%s = ?", COMMENT_ID), secondComment.ParentID).
			Update(REPLIES_COUNT, gorm.Expr("replies_count - 1")).
			Error; err != nil {
			tx.Rollback()
			zlog.Errorf("更新一级评论回复数失败: %v", err)
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		zlog.Errorf("删除评论事务提交失败: %v", err)
		return err
	}

	return nil
}

// GetCommentByID 根据评论ID获取评论信息
func (r *CommentRepo) GetCommentByID(commentID int64) (first *model.FirstComment, second *model.SecondComment, err error) {
	first = &model.FirstComment{}
	second = &model.SecondComment{}

	// 尝试从一级评论表中获取
	if err := r.DB.Model(&model.FirstComment{}).
		Where(fmt.Sprintf("%s = ?", COMMENT_ID), commentID).
		First(first).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果在一级评论表中未找到，尝试从二级评论表中获取
			if err := r.DB.Model(&model.SecondComment{}).
				Where(fmt.Sprintf("%s = ?", COMMENT_ID), commentID).
				First(second).
				Error; err != nil {
				zlog.Errorf("根据ID获取评论失败: %v", err)
				return nil, nil, err
			}
			return nil, second, nil
		}
		zlog.Errorf("根据ID获取一级评论失败: %v", err)
		return nil, nil, err
	}
	return first, nil, nil
}

// GetFirstCommentList 获取一级评论列表
func (r *CommentRepo) GetFirstCommentList(blogID int64, limit int) ([]types.CommentDetail, error) {
	var comments []types.CommentDetail

	err := r.DB.Table("first_comment").
		Select("first_comment.*, user_display.nickname, user_display.avatar, user_display.tag").
		Joins("LEFT JOIN user_display ON first_comment.user_id = user_display.id").
		Where(fmt.Sprintf("first_comment.%s = ? AND first_comment.deleted_at IS NULL", BLOG_ID), blogID).
		Order("first_comment.created_at DESC").
		Limit(limit).
		Scan(&comments).Error

	if err != nil {
		zlog.Errorf("获取一级评论列表失败: %v", err)
		return nil, err
	}

	for i := range comments {
		var secondComments []types.SecondCommentDetail
		err := r.DB.Table("second_comment").
			Select(`
				second_comment.id,
				second_comment.blog_id,
				second_comment.user_id,
				second_comment.content,
				second_comment.parent_id,
				second_comment.created_at,
				user_display.nickname,
				user_display.avatar,
				user_display.tag
			`).
			Joins("LEFT JOIN user_display ON second_comment.user_id = user_display.id").
			Where(fmt.Sprintf("second_comment.%s = ? AND second_comment.deleted_at IS NULL", PARENT_ID), comments[i].ID).
			Order("second_comment.created_at ASC").
			Limit(3).
			Scan(&secondComments).Error

		if err != nil {
			zlog.Errorf("获取二级评论失败: %v", err)
			continue
		}
		comments[i].SecondComments = secondComments
	}

	return comments, nil
}

// GetSecondCommentList 获取更多二级评论
func (r *CommentRepo) GetSecondCommentList(req types.GetSecondCommentListReq) ([]types.SecondCommentDetail, int64, error) {
	var (
		comments   []types.SecondCommentDetail
		totalCount int64
	)

	// 查询二级评论总数
	if err := r.DB.Model(&model.SecondComment{}).
		Where(fmt.Sprintf("%s = ? AND deleted_at IS NULL", PARENT_ID), req.ParentID).
		Count(&totalCount).Error; err != nil {
		zlog.Errorf("获取二级评论总数失败: %v", err)
		return nil, 0, err
	}

	query := r.DB.Table("second_comment").
		Select(`
			second_comment.id,
			second_comment.blog_id,
			second_comment.user_id,
			second_comment.content,
			second_comment.parent_id,
			second_comment.created_at,
			user_display.nickname,
			user_display.avatar,
			user_display.tag
		`).
		Joins("LEFT JOIN user_display ON second_comment.user_id = user_display.id").
		Where(fmt.Sprintf("second_comment.%s = ? AND second_comment.deleted_at IS NULL", PARENT_ID), req.ParentID).
		Order("second_comment.created_at ASC").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize)

	if err := query.Scan(&comments).Error; err != nil {
		zlog.Errorf("获取二级评论列表失败: %v", err)
		return nil, 0, err
	}

	return comments, totalCount, nil
}
func (r *CommentRepo) LikeComment(userID, commentID int64) (*types.LikeCommentResp, error) {
	var isLiked bool = false

	tx := r.DB.Begin()

	// 查询点赞状态
	err := tx.Model(&model.CommentLike{}).
		Select(IS_LIKED).
		Where(fmt.Sprintf("%s = ? AND %s = ?", USER_ID, COMMENT_ID), userID, commentID).
		First(&isLiked).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 初始化点赞记录
			err = tx.Model(&model.CommentLike{}).
				Create(&model.CommentLike{
					UserID:    userID,
					CommentID: commentID,
					IsLiked:   false,
				}).Error
			if err != nil {
				tx.Rollback()
				zlog.Errorf("初始化评论点赞表失败：%v", err)
				return nil, err
			}
		} else {
			tx.Rollback()
			zlog.Errorf("查询评论点赞状态失败：%v", err)
			return nil, err
		}
	}

	// 查询当前点赞数
	var likeCount int64
	err = tx.Model(&model.FirstComment{}).
		Select(LIKE_COUNT).
		Where("id = ?", commentID).
		First(&likeCount).
		Error
	if err != nil {
		tx.Rollback()
		zlog.Errorf("查询评论点赞数失败：%v", err)
		return nil, err
	}

	// 切换点赞状态
	if !isLiked {
		// 执行点赞
		err = tx.Model(&model.CommentLike{}).
			Where(fmt.Sprintf("%s = ? AND %s = ?", USER_ID, COMMENT_ID), userID, commentID).
			Update(IS_LIKED, true).
			Error
		if err != nil {
			tx.Rollback()
			zlog.Errorf("评论点赞操作失败：%v", err)
			return nil, err
		}
		likeCount++
	} else {
		// 取消点赞
		err = tx.Model(&model.CommentLike{}).
			Where(fmt.Sprintf("%s = ? AND %s = ?", USER_ID, COMMENT_ID), userID, commentID).
			Update(IS_LIKED, false).
			Error
		if err != nil {
			tx.Rollback()
			zlog.Errorf("取消评论点赞失败：%v", err)
			return nil, err
		}
		likeCount--
	}

	// 更新评论点赞数
	err = tx.Model(&model.FirstComment{}).
		Where(fmt.Sprintf("%s = ?", COMMENT_ID), commentID).
		Update(LIKE_COUNT, likeCount).
		Error
	if err != nil {
		tx.Rollback()
		zlog.Errorf("更新评论点赞数失败：%v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &types.LikeCommentResp{
		IsLiked:   !isLiked, // 返回操作后的新状态
		LikeCount: likeCount,
	}, nil
}
