package logic

import (
	"context"
	"errors"
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
	codeUnlikeFailed       = response.MsgCode{Code: 40042, Msg: "取消点赞失败"}
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
		Content:        req.Content,
		Tag:            req.Tag,
		SubTag:         req.SubTag,
		ViewPermission: req.ViewPermission,
	}
	//创建帖子
	err = repo.NewBlogRepo(global.DB).CreateBlog(blog)
	if err != nil {
		zlog.CtxErrorf(ctx, "create blog error: %v", err)
		return nil, response.ErrResp(err, codeBlogCreateFailed)
	}
	//存储图片到数据库
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
		Blog:   *blog,
		Images: images,
	}
	return resp, nil
}

// UpdateBlog 更新帖子
func (l *BlogLogic) UpdateBlog(ctx context.Context, req types.UpdateBlogReq, imageUrls []string) (resp *types.UpdateBlogResp, err error) {
	defer utils.RecordTime(time.Now())()

	blog, err := repo.NewBlogRepo(global.DB).GetBlogByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "Blog not found: %v", err)
			return nil, response.ErrResp(err, codeBlogNotFound)
		}
		zlog.CtxErrorf(ctx, "GetBlogByID failed: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}

	blog.Content = req.Content
	blog.Tag = req.Tag
	blog.SubTag = req.SubTag
	blog.ViewPermission = req.ViewPermission

	err = repo.NewBlogRepo(global.DB).UpdateBlog(blog)
	if err != nil {
		zlog.CtxErrorf(ctx, "UpdateBlog failed: %v", err)
		return nil, response.ErrResp(err, codeBlogUpdateFailed)
	}

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

	resp = &types.UpdateBlogResp{
		Blog:   *blog,
		Images: images,
	}
	return resp, nil
}

// DeleteBlog 删除帖子
func (l *BlogLogic) DeleteBlog(ctx context.Context, req types.DeleteBlogReq) (resp *types.DeleteBlogResp, err error) {
	defer utils.RecordTime(time.Now())()

	_, err = repo.NewBlogRepo(global.DB).GetBlogByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "Blog not found: %v", err)
			return nil, response.ErrResp(err, codeBlogNotFound)
		}
		zlog.CtxErrorf(ctx, "GetBlogByID failed: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}

	err = repo.NewBlogRepo(global.DB).DeleteBlog(req.ID)
	if err != nil {
		zlog.CtxErrorf(ctx, "DeleteBlog failed: %v", err)
		return nil, response.ErrResp(err, codeBlogDeleteFailed)
	}

	return resp, nil
}

// GetBlogByID 获取帖子详情
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

	// 获取与帖子关联的所有图片
	images, err := repo.NewImageRepo(global.DB).GetImagesByBlogID(blog.ID)
	if err != nil {
		zlog.CtxErrorf(ctx, "Get images by blog ID failed: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}

	// 构造响应体
	resp = &types.GetBlogByIDResp{
		Blog:   *blog,
		Images: images, // 返回图片列表
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
		Blogs:    blogs,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	return resp, nil
}

// GetBlogsByTag 根据标签获取帖子列表
func (l *BlogLogic) GetBlogsByTag(ctx context.Context, req types.GetBlogsByTagReq) (resp *types.GetBlogsByTagResp, err error) {
	defer utils.RecordTime(time.Now())()

	blogs, err := repo.NewBlogRepo(global.DB).GetBlogsByTag(req.SubTag)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetBlogsByTag failed: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}

	resp = &types.GetBlogsByTagResp{
		Blogs: blogs,
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
		Blogs:    blogs,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	return resp, nil
}

// GetMyBlogs 获取其他用户发布的帖子
func (l *BlogLogic) GetBlogsByUserID(ctx context.Context, req types.GetBlogsByUserIDReq) (resp *types.GetBlogsByUserIDResp, err error) {
	defer utils.RecordTime(time.Now())()

	blogs, total, err := repo.NewBlogRepo(global.DB).GetBlogsByUserID(req.UserID, req.Page, req.PageSize)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetMyBlogs failed: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}

	resp = &types.GetBlogsByUserIDResp{
		Blogs:    blogs,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
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
	resp = &types.CollectBlogResp{
		Success: true,
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
	resp = &types.UncollectBlogResp{
		Success: true,
	}
	return resp, nil
}

// GetCollectedBlogs 获取用户收藏的帖子
func (l *BlogLogic) GetCollectedBlogs(ctx context.Context, req types.GetCollectedBlogsReq, UserID int64) (resp *types.GetCollectedBlogsResp, err error) {
	defer utils.RecordTime(time.Now())()

	blogs, err := repo.NewBlogRepo(global.DB).GetCollectedBlogs(UserID)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetCollectedBlogs failed: %v", err)
		return nil, response.ErrResp(err, codeGetCollectedFailed)
	}

	resp = &types.GetCollectedBlogsResp{
		Blogs: blogs,
	}
	return resp, nil
}

// LikeBlog 点赞帖子
func (l *BlogLogic) LikeBlog(ctx context.Context, req types.LikeBlogReq, UserID int64) (resp *types.LikeBlogResp, err error) {
	defer utils.RecordTime(time.Now())()

	err = repo.NewBlogRepo(global.DB).LikeBlog(UserID, req.BlogID)
	if err != nil {
		zlog.CtxErrorf(ctx, "LikeBlog failed: %v", err)
		return nil, response.ErrResp(err, codeLikeFailed)
	}

	resp = &types.LikeBlogResp{
		Success: true,
	}
	return resp, nil
}

// UnlikeBlog 取消点赞帖子
func (l *BlogLogic) UnlikeBlog(ctx context.Context, req types.UnlikeBlogReq, UserID int64) (resp *types.UnlikeBlogResp, err error) {
	defer utils.RecordTime(time.Now())()

	err = repo.NewBlogRepo(global.DB).UnlikeBlog(UserID, req.BlogID)
	if err != nil {
		zlog.CtxErrorf(ctx, "UnlikeBlog failed: %v", err)
		return nil, response.ErrResp(err, codeUnlikeFailed)
	}

	resp = &types.UnlikeBlogResp{
		Success: true,
	}
	return resp, nil
}
