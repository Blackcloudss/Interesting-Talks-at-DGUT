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
			Update(COMMENT_COUNT, gorm.Expr("%s + 1", COMMENT_COUNT)).
			Error; err != nil {
			tx.Rollback()
			zlog.Errorf("更新博客评论数失败: %v", err)
			return errors.Wrap(err, "更新博客评论数失败")
		}
	} else {
		// 回复评论：更新父评论的 replies_count
		if err := tx.Model(&model.Comment{}).
			Where(fmt.Sprintf("%s = ? AND %s IS NULL", ID, DELETED_AT), comment.ParentID).
			Update(REPLIES_COUNT, gorm.Expr("%s + 1", REPLIES_COUNT)).
			Error; err != nil {
			tx.Rollback()
			zlog.Errorf("更新父评论回复数失败: %v", err)
			return errors.Wrap(err, "更新父评论回复数失败")
		}
	}

	return tx.Commit().Error
}

// GetCommentByID 获取评论基础信息
func (r *CommentRepo) GetCommentByID(id int64) (*model.Comment, error) {
	var comment model.Comment
	err := r.DB.
		Model(&model.Comment{}).
		Where(fmt.Sprintf("comment.%s = ? AND comment.%s IS NULL", ID, DELETED_AT), id).
		First(&comment).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.Warnf("评论不存在 ID:%d", id)
			return nil, errors.New("评论不存在")
		}
		zlog.Errorf("获取评论失败 ID:%d 错误:%v", id, err)
		return nil, errors.Wrap(err, "数据库查询失败")
	}
	return &comment, nil
}

// DeleteComment 删除评论
func (r *CommentRepo) DeleteComment(comment *model.Comment) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		zlog.Errorf("事务启动失败 评论ID:%d 错误:%v", comment.ID, tx.Error)
		return errors.New("事务处理失败")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 执行软删除
	if err := tx.Model(comment).
		Update(DELETED_AT, gorm.Expr("NOW()")).
		Error; err != nil {
		tx.Rollback()
		zlog.Errorf("软删除失败 评论ID:%d 错误:%v", comment.ID, err)
		return errors.New("评论删除失败")
	}

	// 2. 更新计数器
	var updateErr error
	if comment.ParentID == 0 {
		// 主评论：减少博客评论数
		updateErr = tx.Model(&model.Blog{}).
			Where(fmt.Sprintf("blog.%s = ?", ID), comment.BlogID).
			Update(COMMENT_COUNT, gorm.Expr("blog.%s - 1", COMMENT_COUNT)).
			Error
	} else {
		// 子评论：减少父评论回复数
		updateErr = tx.Model(&model.Comment{}).
			Where(fmt.Sprintf("comment.%s = ? AND %s IS NULL", ID, DELETED_AT), comment.ParentID).
			Update(REPLIES_COUNT, gorm.Expr("comment.%s - 1", REPLIES_COUNT)).
			Error
	}

	if updateErr != nil {
		tx.Rollback()
		zlog.Errorf("计数器更新失败 评论ID:%d 错误:%v", comment.ID, updateErr)
		return errors.New("计数器更新失败")
	}

	if err := tx.Commit().Error; err != nil {
		zlog.Errorf("事务提交失败 评论ID:%d 错误:%v", comment.ID, err)
		return errors.New("事务提交失败")
	}

	return nil
}

// GetCommentList 获取顶级评论列表（带分页和排序）
func (r *CommentRepo) GetCommentList(blogID int64, page, pageSize int, sortBy string) ([]types.CommentDetail, int64, error) {
	var comments []types.CommentDetail
	var total int64

	query := r.DB.Model(&model.Comment{}).
		Select(fmt.Sprintf(`
            comment.%s,
            comment.%s,
            comment.%s,
            comment.%s,
            comment.%s,
            comment.%s,
            comment.%s,
            user_display.nickname,
            user_display.avatar,
            user_display.tag
        `, ID, BLOG_ID, USER_ID, CONTENT, LIKE_COUNT, REPLIES_COUNT, CREATED_AT)).
		Joins("LEFT JOIN user_display ON comment.user_id = user_display.id").
		Where(fmt.Sprintf("comment.%s = ? AND comment.%s = 0 AND comment.%s IS NULL", BLOG_ID, PARENT_ID, DELETED_AT), blogID)

	// 添加排序
	switch sortBy {
	case LIKE_COUNT:
		query = query.Order(fmt.Sprintf("comment.%s DESC", LIKE_COUNT))
	default:
		query = query.Order(fmt.Sprintf("comment.%s DESC", CREATED_AT))
	}

	if err := query.Count(&total).Error; err != nil {
		zlog.Errorf("统计评论总数失败: %v", err)
		return nil, 0, err
	}

	if err := query.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&comments).
		Error; err != nil {
		zlog.Errorf("获取评论列表失败: %v", err)
		return nil, 0, err
	}

	return comments, total, nil
}

// GetReplies 获取评论的回复列表（带分页）
func (r *CommentRepo) GetReplies(parentID int64, page, pageSize int) ([]types.ReplyDetail, int64, error) {
	var replies []types.ReplyDetail
	var total int64

	query := r.DB.Model(&model.Comment{}).
		Select(fmt.Sprintf(`
            comment.%s,
            comment.%s,
            comment.%s,
            comment.%s,
            comment.%s,
            comment.%s,
            reply_user.nickname,
            reply_user.avatar,
            reply_user.tag,
            parent.%s AS parent_user_id,
            parent_user.nickname AS parent_nick
        `, ID, BLOG_ID, USER_ID, CONTENT, LIKE_COUNT, CREATED_AT, USER_ID)).
		Joins("LEFT JOIN user_display AS reply_user ON comment.user_id = reply_user.id").
		Joins("LEFT JOIN comment AS parent ON comment.parent_id = parent.id").
		Joins("LEFT JOIN user_display AS parent_user ON parent.user_id = parent_user.id").
		Where(fmt.Sprintf("comment.%s = ? AND comment.%s IS NULL", PARENT_ID, DELETED_AT), parentID).
		Order(fmt.Sprintf("comment.%s ASC", CREATED_AT))

	if err := query.Count(&total).Error; err != nil {
		zlog.Errorf("统计回复总数失败, parentID:%d, error:%v", parentID, err)
		return nil, 0, err
	}

	if err := query.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&replies).
		Error; err != nil {
		zlog.Errorf("获取回复列表失败, parentID:%d, error:%v", parentID, err)
		return nil, 0, err
	}

	return replies, total, nil
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
