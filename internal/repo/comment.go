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

// 数据库列名常量
const (
	PARENT_ID     = "parent_id"     // 父评论ID
	REPLIES_COUNT = "replies_count" // 回复数

)

func NewCommentRepo(db *gorm.DB) *CommentRepo {
	return &CommentRepo{
		DB: db,
	}
}

func (r *CommentRepo) CreateComment(comment *model.Comment) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		zlog.Errorf("创建评论事务开始失败: %v", tx.Error)
		return tx.Error
	}

	// 敏感词过滤
	comment.Content = global.Filter.Replace(comment.Content, '*')

	// 创建评论记录
	if err := tx.Create(comment).Error; err != nil {
		tx.Rollback()
		zlog.Errorf("创建评论记录失败: %v", err)
		return errors.Wrap(err, "创建评论记录失败")
	}

	// 更新计数逻辑
	if comment.ParentID == 0 {
		// 顶级评论：更新博客的 comment_count
		if err := tx.Model(&model.Blog{}).
			Where(fmt.Sprintf("%s = ?", ID), comment.BlogID).
			Update("comment_count", gorm.Expr("comment_count + 1")).
			Error; err != nil {
			tx.Rollback()
			zlog.Errorf("更新博客评论数失败: %v", err)
			return errors.Wrap(err, "更新博客评论数失败")
		}
	} else {
		// 回复评论：更新父评论的 replies_count
		if err := tx.Model(&model.Comment{}).
			Where(fmt.Sprintf("%s = ? AND %s IS NULL", ID, DELETED_AT), comment.ParentID).
			Update(REPLIES_COUNT, gorm.Expr("replies_count + 1")).
			Error; err != nil {
			tx.Rollback()
			zlog.Errorf("更新父评论回复数失败: %v", err)
			return errors.Wrap(err, "更新父评论回复数失败")
		}
	}

	return tx.Commit().Error
}

// DeleteComment 软删除评论并更新相关计数
func (r *CommentRepo) DeleteComment(commentID, blogID, parentID int64) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		zlog.Errorf("删除评论事务开始失败: %v", tx.Error)
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 软删除评论记录
	if err := tx.Model(&model.Comment{}).
		Where(fmt.Sprintf("%s = ?", ID), commentID).
		Update(DELETED_AT, gorm.Expr("NOW()")).
		Error; err != nil {
		tx.Rollback()
		zlog.Errorf("软删除评论失败: %v", err)
		return errors.Wrap(err, "软删除评论失败")
	}

	// 如果是顶级评论，更新博客的评论数
	if parentID == 0 {
		if err := tx.Model(&model.Blog{}).
			Where(fmt.Sprintf("%s = ?", ID), blogID).
			Update("comment_count", gorm.Expr("comment_count - 1")).
			Error; err != nil {
			tx.Rollback()
			zlog.Errorf("更新帖子评论数失败: %v", err)
			return errors.Wrap(err, "更新帖子评论数失败")
		}
	} else {
		// 如果是回复评论，更新父评论的回复数
		if err := tx.Model(&model.Comment{}).
			Where(fmt.Sprintf("%s = ? AND %s IS NULL", ID, DELETED_AT), parentID).
			Update(REPLIES_COUNT, gorm.Expr("replies_count - 1")).
			Error; err != nil {
			tx.Rollback()
			zlog.Errorf("更新父评论回复数失败: %v", err)
			return errors.Wrap(err, "更新父评论回复数失败")
		}
	}

	if err := tx.Commit().Error; err != nil {
		zlog.Errorf("删除评论事务提交失败: %v", err)
		return errors.Wrap(err, "提交事务失败")
	}

	return nil
}

// GetCommentByID 根据评论ID获取评论信息
func (r *CommentRepo) GetCommentByID(commentID int64) (*model.Comment, error) {
	var comment model.Comment

	if err := r.DB.Model(&model.Comment{}).
		Where(fmt.Sprintf("%s = ? AND %s IS NULL", ID, DELETED_AT), commentID).
		First(&comment).
		Error; err != nil {
		zlog.Errorf("根据ID获取评论失败: %v", err)
		return nil, errors.Wrap(err, "未找到评论")
	}

	return &comment, nil
}

// GetCommentList 获取顶级评论列表
func (r *CommentRepo) GetCommentList(blogID int64, page, pageSize int, sortBy string) ([]types.CommentDetail, int64, error) {
	var (
		comments   []types.CommentDetail
		totalCount int64
	)

	// 查询顶级评论总数
	if err := r.DB.Model(&model.Comment{}).
		Where(fmt.Sprintf("%s = ? AND %s IS NULL AND %s = 0", BLOG_ID, DELETED_AT, PARENT_ID), blogID).
		Count(&totalCount).Error; err != nil {
		zlog.Errorf("获取评论总数失败: %v", err)
		return nil, 0, errors.Wrap(err, "获取评论总数失败")
	}

	// 构建查询
	query := r.DB.Table("comment").
		Select("comment.*, user_display.nickname, user_display.avatar, user_display.tag").
		Joins("LEFT JOIN user_display ON comment.user_id = user_display.id").
		Where(fmt.Sprintf("comment.%s = ? AND comment.%s IS NULL AND comment.%s = 0", BLOG_ID, DELETED_AT, PARENT_ID), blogID)

	// 添加排序
	switch sortBy {
	case LIKE_COUNT:
		query = query.Order(fmt.Sprintf("comment.%s DESC", LIKE_COUNT))
	default:
		query = query.Order(fmt.Sprintf("comment.%s DESC", CREATED_AT))
	}

	// 添加分页
	query = query.Offset((page - 1) * pageSize).Limit(pageSize)

	// 执行查询
	if err := query.Scan(&comments).Error; err != nil {
		zlog.Errorf("获取评论列表失败: %v", err)
		return nil, 0, errors.Wrap(err, "获取评论列表失败")
	}

	// 为每个顶级评论获取前3条回复
	for i := range comments {
		var replies []types.CommentDetail
		err := r.DB.Table("comment").
			Select("comment.*, user_display.nickname, user_display.avatar, user_display.tag").
			Joins("LEFT JOIN user_display ON comment.user_id = user_display.id").
			Where(fmt.Sprintf("comment.%s = ? AND comment.%s IS NULL", PARENT_ID, DELETED_AT), comments[i].ID).
			Order(fmt.Sprintf("comment.%s ASC", CREATED_AT)).
			Limit(3).
			Scan(&replies).Error

		if err != nil {
			zlog.Errorf("获取评论回复失败: %v", err)
			continue
		}
		comments[i].Replies = replies
	}

	return comments, totalCount, nil
}

// GetReplies 获取评论的回复列表
func (r *CommentRepo) GetReplies(parentID int64, page, pageSize int) ([]types.CommentDetail, int64, error) {
	var (
		replies    []types.CommentDetail
		totalCount int64
	)

	// 查询回复总数
	if err := r.DB.Model(&model.Comment{}).
		Where(fmt.Sprintf("%s = ? AND %s IS NULL", PARENT_ID, DELETED_AT), parentID).
		Count(&totalCount).Error; err != nil {
		zlog.Errorf("获取回复总数失败: %v", err)
		return nil, 0, errors.Wrap(err, "获取回复总数失败")
	}

	// 构建查询
	query := r.DB.Table("comment").
		Select("comment.*, user_display.nickname, user_display.avatar, user_display.tag").
		Joins("LEFT JOIN user_display ON comment.user_id = user_display.id").
		Where(fmt.Sprintf("comment.%s = ? AND comment.%s IS NULL", PARENT_ID, DELETED_AT), parentID).
		Order(fmt.Sprintf("comment.%s ASC", CREATED_AT)).
		Offset((page - 1) * pageSize).
		Limit(pageSize)

	// 执行查询
	if err := query.Scan(&replies).Error; err != nil {
		zlog.Errorf("获取回复列表失败: %v", err)
		return nil, 0, errors.Wrap(err, "获取回复列表失败")
	}

	return replies, totalCount, nil
}

// LikeComment 点赞/取消点赞评论
func (r *CommentRepo) LikeComment(userID, commentID int64) (*types.LikeCommentResp, error) {
	var isLiked bool = false

	tx := r.DB.Begin()
	if tx.Error != nil {
		zlog.Errorf("点赞评论事务开始失败: %v", tx.Error)
		return nil, tx.Error
	}

	// 查询点赞状态
	err := tx.Model(&model.CommentLike{}).
		Select(IS_LIKED).
		Where(fmt.Sprintf("%s = ? AND %s = ?", USER_ID, ID), userID, commentID).
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
				return nil, errors.Wrap(err, "初始化评论点赞表失败")
			}
		} else {
			tx.Rollback()
			zlog.Errorf("查询评论点赞状态失败：%v", err)
			return nil, errors.Wrap(err, "查询评论点赞状态失败")
		}
	}

	// 查询当前点赞数
	var likeCount int64
	err = tx.Model(&model.Comment{}).
		Select(LIKE_COUNT).
		Where(fmt.Sprintf("%s = ? AND %s IS NULL", ID, DELETED_AT), commentID).
		First(&likeCount).
		Error
	if err != nil {
		tx.Rollback()
		zlog.Errorf("查询评论点赞数失败：%v", err)
		return nil, errors.Wrap(err, "查询评论点赞数失败")
	}

	// 切换点赞状态
	if !isLiked {
		// 执行点赞
		err = tx.Model(&model.CommentLike{}).
			Where(fmt.Sprintf("%s = ? AND %s = ?", USER_ID, ID), userID, commentID).
			Update(IS_LIKED, true).
			Error
		if err != nil {
			tx.Rollback()
			zlog.Errorf("评论点赞操作失败：%v", err)
			return nil, errors.Wrap(err, "评论点赞操作失败")
		}
		likeCount++
	} else {
		// 取消点赞
		err = tx.Model(&model.CommentLike{}).
			Where(fmt.Sprintf("%s = ? AND %s = ?", USER_ID, ID), userID, commentID).
			Update(IS_LIKED, false).
			Error
		if err != nil {
			tx.Rollback()
			zlog.Errorf("取消评论点赞失败：%v", err)
			return nil, errors.Wrap(err, "取消评论点赞失败")
		}
		likeCount--
	}

	// 更新评论点赞数
	err = tx.Model(&model.Comment{}).
		Where(fmt.Sprintf("%s = ? AND %s IS NULL", ID, DELETED_AT), commentID).
		Update(LIKE_COUNT, likeCount).
		Error
	if err != nil {
		tx.Rollback()
		zlog.Errorf("更新评论点赞数失败：%v", err)
		return nil, errors.Wrap(err, "更新评论点赞数失败")
	}

	if err := tx.Commit().Error; err != nil {
		zlog.Errorf("点赞评论事务提交失败: %v", err)
		return nil, errors.Wrap(err, "提交事务失败")
	}

	return &types.LikeCommentResp{
		IsLiked:   !isLiked, // 返回操作后的新状态
		LikeCount: likeCount,
	}, nil
}
