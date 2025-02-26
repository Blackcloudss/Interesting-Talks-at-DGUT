package logic

import (
	"context"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/repo"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// 定义内部逻辑错误
var (
	codeCommentNotFound     = response.MsgCode{Code: 40031, Msg: "评论不存在"}
	codeCommentCreateFailed = response.MsgCode{Code: 40032, Msg: "创建评论失败"}
	codeCommentDeleteFailed = response.MsgCode{Code: 40033, Msg: "删除评论失败"}
	codeUnauthorized        = response.MsgCode{Code: 40034, Msg: "用户未授权"}
	codeUserNotLoggedIn     = response.MsgCode{Code: 20001, Msg: "用户未登录"}
)

// CreateComment 创建评论
func CreateComment(ctx context.Context, comment *model.Comment) error {
	tx := global.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		zlog.CtxErrorf(ctx, "Failed to start transaction: %v", tx.Error)
		return response.ErrResp(tx.Error, codeCommentCreateFailed)
	}

	// 检查用户是否登录
	if comment.AuthorID == 0 {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "CreateComment failed: user not logged in")
		return response.ErrResp(nil, codeUserNotLoggedIn)
	}

	// 检查帖子是否存在
	if _, err := repo.GetBlogByID(comment.BlogID); err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "CreateComment failed: blog not found (blogID: %d)", comment.BlogID)
			return response.ErrResp(err, codeBlogNotFound)
		}
		zlog.CtxErrorf(ctx, "CreateComment failed: %v", err)
		return response.ErrResp(err, codeCommentCreateFailed)
	}

	// 使用雪花算法生成评论ID
	commentID := global.Node.Generate().Int64()
	comment.CommentID = uint64(commentID)

	// 创建评论
	if err := repo.CreateComment(tx, comment); err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "CreateComment failed: %v", err)
		return response.ErrResp(err, codeCommentCreateFailed)
	}

	// 更新帖子的评论数
	if err := repo.IncrementCommentCount(tx, comment.BlogID); err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "IncrementCommentCount failed: %v", err)
		return response.ErrResp(err, codeCommentCreateFailed)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "Failed to commit transaction: %v", err)
		return response.ErrResp(err, codeCommentCreateFailed)
	}

	zlog.CtxInfof(ctx, "Comment created successfully (commentID: %d, blogID: %d)", comment.CommentID, comment.BlogID)
	return nil
}

// DeleteComment 删除评论
func DeleteComment(ctx context.Context, commentID uint64, userID uint64) error {
	tx := global.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		zlog.CtxErrorf(ctx, "Failed to start transaction: %v", tx.Error)
		return response.ErrResp(tx.Error, codeCommentDeleteFailed)
	}

	// 获取评论信息
	comment, err := repo.GetCommentByID(ctx, commentID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "DeleteComment failed: comment not found (commentID: %d)", commentID)
			return response.ErrResp(err, codeCommentNotFound)
		}
		zlog.CtxErrorf(ctx, "GetCommentByID failed: %v", err)
		return response.ErrResp(err, codeCommentDeleteFailed)
	}

	// 检查用户是否是评论作者或管理员
	if comment.AuthorID != userID {
		userRole, err := repo.GetUserRole(userID)
		if err != nil {
			tx.Rollback()
			zlog.CtxErrorf(ctx, "GetUserRole failed: %v", err)
			return response.ErrResp(err, codeUnauthorized)
		}
		if userRole != "管理员" {
			tx.Rollback()
			zlog.CtxErrorf(ctx, "DeleteComment failed: user not authorized (userID: %d, commentID: %d)", userID, commentID)
			return response.ErrResp(nil, codeUnauthorized)
		}
	}

	// 删除评论
	if err := repo.DeleteComment(tx, commentID); err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "DeleteComment failed: %v", err)
		return response.ErrResp(err, codeCommentDeleteFailed)
	}

	// 更新帖子的评论数
	if err := repo.DecrementCommentCount(tx, comment.BlogID); err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "DecrementCommentCount failed: %v", err)
		return response.ErrResp(err, codeCommentDeleteFailed)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "Failed to commit transaction: %v", err)
		return response.ErrResp(err, codeCommentDeleteFailed)
	}

	zlog.CtxInfof(ctx, "Comment deleted successfully (commentID: %d, blogID: %d)", commentID, comment.BlogID)
	return nil
}

// GetCommentList 获取评论列表
func GetCommentList(ctx context.Context, blogID uint64) ([]model.Comment, error) {
	comments, err := repo.GetCommentList(ctx, blogID)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetCommentList failed: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}

	zlog.CtxInfof(ctx, "Comment list retrieved successfully (blogID: %d, count: %d)", blogID, len(comments))
	return comments, nil
}
