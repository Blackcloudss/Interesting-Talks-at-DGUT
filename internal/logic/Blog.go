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

// BlogLogic 帖子逻辑层
type BlogLogic struct{}

// NewBlogLogic 创建帖子逻辑层实例
func NewBlogLogic() *BlogLogic {
	return &BlogLogic{}
}

// 定义内部逻辑错误
var (
	codeBlogCreateFailed = response.MsgCode{Code: 40021, Msg: "创建帖子失败"}
	codeBlogUpdateFailed = response.MsgCode{Code: 40022, Msg: "更新帖子失败"}
	codeBlogDeleteFailed = response.MsgCode{Code: 40023, Msg: "删除帖子失败"}
)

// CreateBlog 创建帖子
func (l *BlogLogic) CreateBlog(ctx context.Context, blog *model.Blog) error {
	tx := global.DB.Begin()
	if tx.Error != nil {
		zlog.CtxErrorf(ctx, "Failed to start transaction: %v", tx.Error)
		return response.ErrResp(tx.Error, codeTransactionFailed)
	}

	// 生成帖子ID（雪花算法）
	blogID := global.Node.Generate().Int64()
	blog.ID = uint64(blogID)

	// 创建帖子
	if err := repo.CreateBlog(tx, blog); err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "Failed to create blog: %v", err)
		return response.ErrResp(err, codeBlogCreateFailed)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		zlog.CtxErrorf(ctx, "Failed to commit transaction: %v", err)
		return response.ErrResp(err, codeTransactionFailed)
	}

	zlog.CtxInfof(ctx, "Blog created successfully: %+v", blog)
	return nil
}

// UpdateBlog 更新帖子
func (l *BlogLogic) UpdateBlog(ctx context.Context, blog *model.Blog) error {
	// 检查帖子是否存在
	existingBlog, err := repo.GetBlogByID(blog.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "Blog not found: %v", err)
			return response.ErrResp(err, codeBlogNotFound)
		}
		zlog.CtxErrorf(ctx, "repo.GetBlogByID failed: %v", err)
		return response.ErrResp(err, response.INTERNAL_ERROR)
	}

	// 更新帖子内容
	existingBlog.Content = blog.Content
	existingBlog.Tag = blog.Tag
	existingBlog.SubTag = blog.SubTag
	existingBlog.ViewPermission = blog.ViewPermission

	// 调用数据库操作更新帖子
	if err := repo.UpdateBlog(existingBlog); err != nil {
		zlog.CtxErrorf(ctx, "repo.UpdateBlog failed: %v", err)
		return response.ErrResp(err, codeBlogUpdateFailed)
	}

	zlog.CtxInfof(ctx, "Blog updated successfully: %+v", existingBlog)
	return nil
}

// DeleteBlog 删除帖子
func (l *BlogLogic) DeleteBlog(ctx context.Context, blogID uint64) error {
	// 获取帖子详情
	blog, err := repo.GetBlogByID(blogID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "Blog not found: %v", err)
			return response.ErrResp(err, codeBlogNotFound)
		}
		zlog.CtxErrorf(ctx, "repo.GetBlogByID failed: %v", err)
		return response.ErrResp(err, response.INTERNAL_ERROR)
	}

	// 删除帖子
	if err := repo.DeleteBlog(blogID); err != nil {
		zlog.CtxErrorf(ctx, "repo.DeleteBlog failed: %v", err)
		return response.ErrResp(err, codeBlogDeleteFailed)
	}

	zlog.CtxInfof(ctx, "Blog deleted successfully: %+v", blog)
	return nil
}

// GetBlogByID 获取帖子详情
func (l *BlogLogic) GetBlogByID(ctx context.Context, blogID uint64, userID uint64) (*model.Blog, error) {
	blog, err := repo.GetBlogByID(blogID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "Blog not found: %v", err)
			return nil, response.ErrResp(err, codeBlogNotFound)
		}
		zlog.CtxErrorf(ctx, "repo.GetBlogByID failed: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}

	// 检查观看权限
	if err := checkViewPermission(ctx, blog, userID); err != nil {
		zlog.CtxWarnf(ctx, "User does not have permission to view this blog: %+v", blog)
		return nil, response.ErrResp(err, codeInsufficientPermission)
	}

	zlog.CtxInfof(ctx, "Blog retrieved successfully: %+v", blog)
	return blog, nil
}

// GetBlogs 分页获取帖子列表
func (l *BlogLogic) GetBlogs(ctx context.Context, page, pageSize int, userID uint64) ([]model.Blog, int64, error) {
	blogs, total, err := repo.GetBlogs(page, pageSize)
	if err != nil {
		zlog.CtxErrorf(ctx, "repo.GetBlogs failed: %v", err)
		return nil, 0, response.ErrResp(err, response.INTERNAL_ERROR)
	}

	// 过滤帖子，根据观看权限
	filteredBlogs := make([]model.Blog, 0)
	for _, blog := range blogs {
		if err := checkViewPermission(ctx, &blog, userID); err == nil {
			filteredBlogs = append(filteredBlogs, blog)
		}
	}

	zlog.CtxInfof(ctx, "Blogs retrieved successfully: %+v", filteredBlogs)
	return filteredBlogs, total, nil
}

// GetBlogsAfterID 获取比指定 ID 更新的帖子
func (l *BlogLogic) GetBlogsAfterID(ctx context.Context, latestID uint64, pageSize int) ([]model.Blog, error) {
	blogs, err := repo.GetBlogsAfterID(latestID, pageSize)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get blogs after ID: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}
	return blogs, nil
}

// GetBlogsByTag 根据标签分页获取帖子列表
func (l *BlogLogic) GetBlogsByTag(ctx context.Context, tag string, page, pageSize int) ([]model.Blog, int64, error) {
	blogs, total, err := repo.GetBlogsByTag(tag, page, pageSize)
	if err != nil {
		zlog.CtxErrorf(ctx, "repo.GetBlogsByTag failed: %v", err)
		return nil, 0, response.ErrResp(err, response.INTERNAL_ERROR)
	}
	zlog.CtxInfof(ctx, "Blogs retrieved successfully: %+v", blogs)
	return blogs, total, nil
}

// GetBlogsByUserID 根据用户ID分页获取帖子列表
func (l *BlogLogic) GetBlogsByUserID(ctx context.Context, userID uint64, page, pageSize int) ([]model.Blog, int64, error) {
	blogs, total, err := repo.GetBlogsByUserID(userID, page, pageSize)
	if err != nil {
		zlog.CtxErrorf(ctx, "repo.GetBlogsByUserID failed: %v", err)
		return nil, 0, response.ErrResp(err, response.INTERNAL_ERROR)
	}
	zlog.CtxInfof(ctx, "Blogs retrieved successfully: %+v", blogs)
	return blogs, total, nil
}

// checkViewPermission 检查用户是否有权限查看帖子
func checkViewPermission(ctx context.Context, blog *model.Blog, userID uint64) error {
	switch blog.ViewPermission {
	case "所有人":
		return nil
	case "仅学生":
		if userID == 0 {
			// 如果是游客（未登录用户），直接拒绝访问
			return response.ErrResp(nil, codeInsufficientPermission)
		}
		role, err := repo.GetUserRole(userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 用户不存在，可能是游客或非法用户
				zlog.CtxWarnf(ctx, "User not found (userID: %d)", userID)
				return response.ErrResp(err, codeUserNotFound)
			}
			zlog.CtxErrorf(ctx, "Failed to get user role: %v", err)
			return response.ErrResp(err, response.INTERNAL_ERROR)
		}
		if role != "student" {
			return response.ErrResp(nil, codeInsufficientPermission)
		}
	case "仅好友":
		if userID == blog.UserID {
			return nil // 自己可以查看自己的帖子
		}
		isFriend := repo.IsFollowing(userID, blog.UserID) && repo.IsFollowing(blog.UserID, userID)
		if !isFriend {
			return response.ErrResp(nil, codeInsufficientPermission)
		}
	case "仅自己":
		if userID != blog.UserID {
			return response.ErrResp(nil, codeInsufficientPermission)
		}
	default:
		return response.ErrResp(nil, response.INTERNAL_ERROR)
	}
	return nil
}
