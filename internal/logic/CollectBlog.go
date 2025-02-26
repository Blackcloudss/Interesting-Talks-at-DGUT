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
	codeCollectFailed      = response.MsgCode{Code: 40026, Msg: "收藏帖子失败"}
	codeUncollectFailed    = response.MsgCode{Code: 40027, Msg: "取消收藏失败"}
	codeGetCollectedFailed = response.MsgCode{Code: 40028, Msg: "获取收藏帖子失败"}
	codeAlreadyCollected   = response.MsgCode{Code: 40029, Msg: "帖子已收藏"}
	codeNotCollected       = response.MsgCode{Code: 40030, Msg: "帖子未收藏"}
)

// CollectBlog 收藏帖子
func CollectBlog(ctx context.Context, userID, blogID uint64) error {
	tx := global.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		zlog.CtxErrorf(ctx, "Failed to start transaction: %v", tx.Error)
		return response.ErrResp(tx.Error, codeCollectFailed)
	}

	// 检查帖子是否存在
	blog, err := repo.GetBlogByID(blogID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "CollectBlog failed: blog not found (blogID: %d)", blogID)
			return response.ErrResp(err, codeBlogNotFound)
		}
		zlog.CtxErrorf(ctx, "CollectBlog failed: %v", err)
		return response.ErrResp(err, codeCollectFailed)
	}

	// 检查是否已经收藏
	if repo.IsCollected(tx, userID, blogID) {
		tx.Rollback()
		zlog.CtxInfof(ctx, "CollectBlog skipped: already collected (userID: %d, blogID: %d)", userID, blogID)
		return response.ErrResp(nil, codeAlreadyCollected)
	}

	// 插入收藏记录
	if err := repo.CollectBlog(tx, userID, blogID); err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "CollectBlog failed: %v", err)
		return response.ErrResp(err, codeCollectFailed)
	}

	// 更新帖子的收藏数
	blog.BeCollected++
	if err := repo.UpdateBlog(blog); err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "UpdateBlog failed: %v", err)
		return response.ErrResp(err, codeCollectFailed)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "Failed to commit transaction: %v", err)
		return response.ErrResp(err, codeCollectFailed)
	}

	zlog.CtxInfof(ctx, "Blog collected successfully (userID: %d, blogID: %d)", userID, blogID)
	return nil
}

// UncollectBlog 取消收藏
func UncollectBlog(ctx context.Context, userID, blogID uint64) error {
	tx := global.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		zlog.CtxErrorf(ctx, "Failed to start transaction: %v", tx.Error)
		return response.ErrResp(tx.Error, codeUncollectFailed)
	}

	// 检查帖子是否存在
	blog, err := repo.GetBlogByID(blogID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "UncollectBlog failed: blog not found (blogID: %d)", blogID)
			return response.ErrResp(err, codeBlogNotFound)
		}
		zlog.CtxErrorf(ctx, "UncollectBlog failed: %v", err)
		return response.ErrResp(err, codeUncollectFailed)
	}

	// 检查是否已经取消收藏
	if !repo.IsCollected(tx, userID, blogID) {
		tx.Rollback()
		zlog.CtxInfof(ctx, "UncollectBlog skipped: not collected (userID: %d, blogID: %d)", userID, blogID)
		return response.ErrResp(nil, codeNotCollected)
	}

	// 执行取消收藏操作
	if err := repo.UncollectBlog(tx, userID, blogID); err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "UncollectBlog failed: %v", err)
		return response.ErrResp(err, codeUncollectFailed)
	}

	// 更新帖子的收藏数
	blog.BeCollected--
	if blog.BeCollected < 0 {
		blog.BeCollected = 0 // 防止收藏数为负
	}
	if err := repo.UpdateBlog(blog); err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "UpdateBlog failed: %v", err)
		return response.ErrResp(err, codeUncollectFailed)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "Failed to commit transaction: %v", err)
		return response.ErrResp(err, codeUncollectFailed)
	}

	zlog.CtxInfof(ctx, "Blog uncollected successfully (userID: %d, blogID: %d)", userID, blogID)
	return nil
}

// GetCollectedBlogs 获取用户收藏的帖子
func GetCollectedBlogs(ctx context.Context, userID uint64) ([]model.Blog, error) {
	// 调用数据层方法获取收藏的帖子
	blogs, err := repo.GetCollectedBlogs(ctx, userID)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get collected blogs for user (userID: %d): %v", userID, err)
		return nil, response.ErrResp(err, codeGetCollectedFailed)
	}

	zlog.CtxInfof(ctx, "Collected blogs retrieved successfully (userID: %d)", userID)
	return blogs, nil
}
