package logic

import (
	"context"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/repo"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// 定义内部逻辑错误
var (
	codeLikeFailed   = response.MsgCode{Code: 40041, Msg: "点赞失败"}
	codeUnlikeFailed = response.MsgCode{Code: 40042, Msg: "取消点赞失败"}
	codeAlreadyLiked = response.MsgCode{Code: 40043, Msg: "已经点赞"}
	codeNotLiked     = response.MsgCode{Code: 40044, Msg: "未点赞"}
)

// LikeBlog 点赞帖子
func LikeBlog(ctx context.Context, userID, blogID uint64) error {
	tx := global.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		zlog.CtxErrorf(ctx, "Failed to start transaction: %v", tx.Error)
		return response.ErrResp(tx.Error, codeLikeFailed)
	}

	// 检查帖子是否存在
	if _, err := repo.GetBlogByID(blogID); err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "LikeBlog failed: blog not found (blogID: %d)", blogID)
			return response.ErrResp(err, codeBlogNotFound)
		}
		zlog.CtxErrorf(ctx, "LikeBlog failed: %v", err)
		return response.ErrResp(err, codeLikeFailed)
	}

	// 检查是否已经点赞
	if repo.IsLiked(tx, userID, blogID) {
		tx.Rollback()
		zlog.CtxInfof(ctx, "LikeBlog skipped: already liked (userID: %d, blogID: %d)", userID, blogID)
		return response.ErrResp(nil, codeAlreadyLiked)
	}

	// 创建点赞记录
	if err := repo.LikeBlog(tx, userID, blogID); err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "LikeBlog failed: %v", err)
		return response.ErrResp(err, codeLikeFailed)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "Failed to commit transaction: %v", err)
		return response.ErrResp(err, codeLikeFailed)
	}

	zlog.CtxInfof(ctx, "Blog liked successfully (userID: %d, blogID: %d)", userID, blogID)
	return nil
}

// UnlikeBlog 取消点赞
func UnlikeBlog(ctx context.Context, userID, blogID uint64) error {
	tx := global.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		zlog.CtxErrorf(ctx, "Failed to start transaction: %v", tx.Error)
		return response.ErrResp(tx.Error, codeUnlikeFailed)
	}

	// 检查帖子是否存在
	if _, err := repo.GetBlogByID(blogID); err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "UnlikeBlog failed: blog not found (blogID: %d)", blogID)
			return response.ErrResp(err, codeBlogNotFound)
		}
		zlog.CtxErrorf(ctx, "UnlikeBlog failed: %v", err)
		return response.ErrResp(err, codeUnlikeFailed)
	}

	// 检查是否已经点赞
	if !repo.IsLiked(tx, userID, blogID) {
		tx.Rollback()
		zlog.CtxInfof(ctx, "UnlikeBlog skipped: not liked (userID: %d, blogID: %d)", userID, blogID)
		return response.ErrResp(nil, codeNotLiked)
	}

	// 执行取消点赞操作
	if err := repo.UnlikeBlog(tx, userID, blogID); err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "UnlikeBlog failed: %v", err)
		return response.ErrResp(err, codeUnlikeFailed)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "Failed to commit transaction: %v", err)
		return response.ErrResp(err, codeUnlikeFailed)
	}

	zlog.CtxInfof(ctx, "Blog unliked successfully (userID: %d, blogID: %d)", userID, blogID)
	return nil
}
