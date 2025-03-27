package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/repo"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils"
	"gorm.io/gorm"
	"time"
)

// 定义内部逻辑错误
var (
	codeBlogCreateFailed   = response.MsgCode{Code: 40021, Msg: "创建帖子失败"}
	codeBlogUpdateFailed   = response.MsgCode{Code: 40022, Msg: "更新帖子失败"}
	codeBlogDeleteFailed   = response.MsgCode{Code: 40023, Msg: "删除帖子失败"}
	codeImageCreateFailed  = response.MsgCode{Code: 40024, Msg: "上传图片失败"}
	codeCollectFailed      = response.MsgCode{Code: 40026, Msg: "收藏帖子失败"}
	codeUncollectFailed    = response.MsgCode{Code: 40027, Msg: "取消收藏失败"}
	codeGetCollectedFailed = response.MsgCode{Code: 40028, Msg: "获取收藏帖子失败"}
	codeLikeFailed         = response.MsgCode{Code: 40041, Msg: "点赞失败"}
)

type BlogLogic struct{}

func NewBlogLogic() *BlogLogic {
	return &BlogLogic{}
}

// CreateBlog 创建帖子
func (l *BlogLogic) CreateBlog(ctx context.Context, req types.CreateBlogReq, UserID int64, imageUrls []string) (resp *types.CreateBlogResp, err error) {
	defer utils.RecordTime(time.Now())()

	blog := &model.Blog{
		UserID:         UserID,
		Title:          req.Title,
		Content:        req.Content,
		BlogTag:        req.BlogTag,
		SubTag:         req.SubTag,
		ViewPermission: req.ViewPermission,
	}
	// 创建帖子
	err = repo.NewBlogRepo(global.DB).CreateBlog(blog)
	if err != nil {
		zlog.CtxErrorf(ctx, "create blog error: %v", err)
		return nil, response.ErrResp(err, codeBlogCreateFailed)
	}

	// 存储图片到数据库
	var images []model.Image
	for _, imageUrl := range imageUrls {
		image := model.Image{
			ImagePath: imageUrl,
			BlogID:    blog.ID,
		}
		err = repo.NewImageRepo(global.DB).CreateImage(&image)
		if err != nil {
			zlog.CtxErrorf(ctx, "create image error: %v", err)
			return nil, response.ErrResp(err, codeImageCreateFailed)
		}
		images = append(images, image)
	}

	resp = &types.CreateBlogResp{
		BlogID:   blog.ID,
		CreateAt: blog.CreatedAt,
	}
	return resp, nil
}

// UpdateBlog 更新帖子
func (l *BlogLogic) UpdateBlog(ctx context.Context, req types.UpdateBlogReq, imageUrls []string) (resp *types.UpdateBlogResp, err error) {
	defer utils.RecordTime(time.Now())()

	// 检查帖子是否存在
	blog, err := repo.NewBlogRepo(global.DB).CheckBlogExists(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "Blog not found: %v", err)
			return nil, response.ErrResp(err, codeBlogNotFound)
		}
		zlog.CtxErrorf(ctx, "CheckBlogExists failed: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}

	// 更新帖子内容
	blog.Title = req.Title
	blog.Content = req.Content
	blog.BlogTag = req.BlogTag
	blog.SubTag = req.SubTag
	blog.ViewPermission = req.ViewPermission

	// 更新帖子记录
	err = repo.NewBlogRepo(global.DB).UpdateBlog(blog)
	if err != nil {
		zlog.CtxErrorf(ctx, "UpdateBlog failed: %v", err)
		return nil, response.ErrResp(err, codeBlogUpdateFailed)
	}

	// 删除旧图片记录
	err = repo.NewImageRepo(global.DB).DeleteImagesByBlogID(blog.ID)
	if err != nil {
		zlog.CtxErrorf(ctx, "Delete images by blog ID failed: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}

	// 创建新图片记录
	for _, imageUrl := range imageUrls {
		image := model.Image{
			ImagePath: imageUrl,
			BlogID:    blog.ID,
		}
		err = repo.NewImageRepo(global.DB).CreateImage(&image)
		if err != nil {
			zlog.CtxErrorf(ctx, "create image error: %v", err)
			return nil, response.ErrResp(err, codeImageCreateFailed)
		}
	}

	// 构造响应体
	resp = &types.UpdateBlogResp{
		UpdateAt: blog.UpdatedAt, // 使用更新时间
	}
	return resp, nil
}

// DeleteBlog 删除帖子
func (l *BlogLogic) DeleteBlog(ctx context.Context, req types.DeleteBlogReq) (resp *types.DeleteBlogResp, err error) {
	defer utils.RecordTime(time.Now())()

	// 检查帖子是否存在
	_, err = repo.NewBlogRepo(global.DB).CheckBlogExists(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "Blog not found: %v", err)
			return nil, response.ErrResp(err, codeBlogNotFound)
		}
		zlog.CtxErrorf(ctx, "CheckBlogExists failed: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}

	// 删除与帖子关联的图片记录
	err = repo.NewImageRepo(global.DB).DeleteImagesByBlogID(req.ID)
	if err != nil {
		zlog.CtxErrorf(ctx, "Delete images by blog ID failed: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}

	// 删除帖子记录
	err = repo.NewBlogRepo(global.DB).DeleteBlog(req.ID)
	if err != nil {
		zlog.CtxErrorf(ctx, "DeleteBlog failed: %v", err)
		return nil, response.ErrResp(err, codeBlogDeleteFailed)
	}

	return resp, nil
}

// GetBlogByID 根据ID获取帖子详情
func (l *BlogLogic) GetBlogByID(ctx context.Context, req types.GetBlogByIDReq) (resp *types.GetBlogByIDResp, err error) {
	defer utils.RecordTime(time.Now())()

	// 获取帖子详情
	blog, err := repo.NewBlogRepo(global.DB).GetBlogByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "Blog not found: %v", err)
			return nil, response.ErrResp(err, codeBlogNotFound)
		}
		zlog.CtxErrorf(ctx, "GetBlogByID failed: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}

	// 获取与帖子关联的所有图片路径
	imagePaths, err := repo.NewImageRepo(global.DB).GetImagePathsByBlogID(req.ID)
	if err != nil {
		zlog.CtxErrorf(ctx, "Get image paths by blog ID failed: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}

	// 构造响应体
	resp = &types.GetBlogByIDResp{
		BlogResp: blog,
		Images:   imagePaths,
	}
	return resp, nil
}

// GetBlogs 分页获取帖子列表
func (l *BlogLogic) GetBlogs(ctx context.Context, req types.GetBlogsReq) (resp *types.GetBlogsResp, err error) {
	defer utils.RecordTime(time.Now())()

	blogs, total, err := repo.NewBlogRepo(global.DB).GetBlogs(req.Page, req.PageSize)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetBlogs failed: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}

	resp = &types.GetBlogsResp{
		Blogs: blogs,
		Total: total,
		PageReq: types.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}
	return resp, nil
}

// GetBlogsByTag 根据标签获取帖子列表
func (l *BlogLogic) GetBlogsByTag(ctx context.Context, req types.GetBlogsByTagReq) (resp *types.GetBlogsByTagResp, err error) {
	defer utils.RecordTime(time.Now())()

	// 获取帖子列表
	blogs, total, err := repo.NewBlogRepo(global.DB).GetBlogsByTag(req.SubTag, req.Page, req.PageSize)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetBlogsByTag failed: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}

	resp = &types.GetBlogsByTagResp{
		Blogs: blogs,
		Total: total,
		PageReq: types.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}

	return resp, nil
}

// GetMyBlogs 获取当前用户发布的帖子
func (l *BlogLogic) GetMyBlogs(ctx context.Context, req types.GetMyBlogsReq, UserID int64) (resp *types.GetMyBlogsResp, err error) {
	defer utils.RecordTime(time.Now())()

	blogs, total, err := repo.NewBlogRepo(global.DB).GetBlogsByUserID(UserID, req.Page, req.PageSize)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetMyBlogs failed: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}

	resp = &types.GetMyBlogsResp{
		Blogs: blogs,
		Total: total,
		PageReq: types.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}
	return resp, nil
}

// GetBlogsByUserID 获取其他用户发布的帖子
func (l *BlogLogic) GetBlogsByUserID(ctx context.Context, req types.GetBlogsByUserIDReq) (resp *types.GetBlogsByUserIDResp, err error) {
	defer utils.RecordTime(time.Now())()

	blogs, total, err := repo.NewBlogRepo(global.DB).GetBlogsByUserID(req.UserID, req.Page, req.PageSize)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetBlogsByUserID failed: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}

	resp = &types.GetBlogsByUserIDResp{
		Blogs: blogs,
		Total: total,
		PageReq: types.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}
	return resp, nil
}

// CollectBlog 收藏帖子
func (l *BlogLogic) CollectBlog(ctx context.Context, req types.CollectBlogReq, UserID int64) (resp *types.CollectBlogResp, err error) {
	defer utils.RecordTime(time.Now())()

	err = repo.NewBlogRepo(global.DB).CollectBlog(UserID, req.BlogID)
	if err != nil {
		zlog.CtxErrorf(ctx, "CollectBlog failed: %v", err)
		return nil, response.ErrResp(err, codeCollectFailed)
	}
	return resp, nil
}

// UncollectBlog 取消收藏帖子
func (l *BlogLogic) UncollectBlog(ctx context.Context, req types.UncollectBlogReq, UserID int64) (resp *types.UncollectBlogResp, err error) {
	defer utils.RecordTime(time.Now())()

	err = repo.NewBlogRepo(global.DB).UncollectBlog(UserID, req.BlogID)
	if err != nil {
		zlog.CtxErrorf(ctx, "UncollectBlog failed:%v", err)
		return nil, response.ErrResp(err, codeUncollectFailed)
	}
	return resp, nil
}

// GetCollectedBlogs 获取用户收藏的帖子
func (l *BlogLogic) GetCollectedBlogs(ctx context.Context, req types.GetCollectedBlogsReq, UserID int64) (resp *types.GetCollectedBlogsResp, err error) {
	defer utils.RecordTime(time.Now())()

	// 获取用户收藏的帖子列表（分页）
	blogs, total, err := repo.NewBlogRepo(global.DB).GetCollectedBlogs(UserID, req.Page, req.PageSize)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetCollectedBlogs failed: %v", err)
		return nil, response.ErrResp(err, codeGetCollectedFailed)
	}

	resp = &types.GetCollectedBlogsResp{
		Blogs: blogs,
		Total: total,
		PageReq: types.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}
	return resp, nil
}

// LikeBlog 点赞或取消点赞帖子
func (l *BlogLogic) LikeBlog(ctx context.Context, BlogID int64, isLiked bool, UserID int64) (resp *types.LikeBlogResp, err error) {
	defer utils.RecordTime(time.Now())()
	// 用 redis 加锁
	lockKey := fmt.Sprintf("blog:like:lock:user:%d:blog:%d", UserID, BlogID)
	locked, err := global.Rdb.SetNX(ctx, lockKey, 1, 1*time.Second).Result()
	if err != nil {
		zlog.CtxErrorf(ctx, "Redis 上锁失败: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}
	if !locked {
		// 未获取到锁，说明该操作正在被其他请求处理
		zlog.CtxInfof(ctx, "点赞/取消赞操作正被 user: %d, blog: %d 使用，请稍等 1 s", UserID, BlogID)
		return nil, response.ErrResp(err, response.USER_OPERATION_LOCKED) // 用户操作被锁定
	}
	defer global.Rdb.Del(ctx, lockKey)

	resp, err = repo.NewBlogRepo(global.DB).LikeBlog(UserID, BlogID)
	if err != nil {
		zlog.CtxErrorf(ctx, "ToggleLikeBlog failed: %v", err)
		return nil, response.ErrResp(err, codeLikeFailed)
	}
	return resp, nil
}
