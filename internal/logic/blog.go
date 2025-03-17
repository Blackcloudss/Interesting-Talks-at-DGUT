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
func (l *BlogLogic) CreateBlog(ctx context.Context, req types.CreateBlogReq) (resp *types.CreateBlogResp, err error) {
	defer utils.RecordTime(time.Now())()

	blog := &model.Blog{
		UserID:         int64(req.UserID),
		Content:        req.Content,
		Tag:            req.Tag,
		SubTag:         req.SubTag,
		ViewPermission: req.ViewPermission,
	}

	err = repo.NewBlogRepo(global.DB).CreateBlog(blog)
	if err != nil {
		zlog.CtxErrorf(ctx, "create blog error: %v", err)
		return nil, response.ErrResp(err, codeBlogCreateFailed)
	}

	resp = &types.CreateBlogResp{
		Blog: *blog,
	}
	return resp, nil
}

// UpdateBlog 更新帖子
func (l *BlogLogic) UpdateBlog(ctx context.Context, req types.UpdateBlogReq) (resp *types.UpdateBlogResp, err error) {
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

	resp = &types.UpdateBlogResp{
		Blog: *blog,
	}
	return resp, nil
}

// DeleteBlog 删除帖子
func (l *BlogLogic) DeleteBlog(ctx context.Context, req types.DeleteBlogReq) (err error) {
	defer utils.RecordTime(time.Now())()

	_, err = repo.NewBlogRepo(global.DB).GetBlogByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "Blog not found: %v", err)
			return response.ErrResp(err, codeBlogNotFound)
		}
		zlog.CtxErrorf(ctx, "GetBlogByID failed: %v", err)
		return response.ErrResp(err, response.INTERNAL_ERROR)
	}

	err = repo.NewBlogRepo(global.DB).DeleteBlog(req.ID)
	if err != nil {
		zlog.CtxErrorf(ctx, "DeleteBlog failed: %v", err)
		return response.ErrResp(err, codeBlogDeleteFailed)
	}

	return nil
}

// GetBlogByID 获取帖子详情
func (l *BlogLogic) GetBlogByID(ctx context.Context, req types.GetBlogByIDReq) (resp *types.GetBlogByIDResp, err error) {
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

	resp = &types.GetBlogByIDResp{
		Blog: *blog,
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
		List:     blogs,
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
		List: blogs,
	}
	return resp, nil
}

// GetMyBlogs 获取当前用户发布的帖子
func (l *BlogLogic) GetMyBlogs(ctx context.Context, req types.GetMyBlogsReq) (resp *types.GetMyBlogsResp, err error) {
	defer utils.RecordTime(time.Now())()

	blogs, total, err := repo.NewBlogRepo(global.DB).GetBlogsByUserID(int64(req.UserID), req.Page, req.PageSize)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetMyBlogs failed: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}

	resp = &types.GetMyBlogsResp{
		List:     blogs,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	return resp, nil
}

// CollectBlog 收藏帖子
func (l *BlogLogic) CollectBlog(ctx context.Context, req types.CollectBlogReq) error {
	defer utils.RecordTime(time.Now())()

	err := repo.NewBlogRepo(global.DB).CollectBlog(int64(req.UserID), req.BlogID)
	if err != nil {
		zlog.CtxErrorf(ctx, "CollectBlog failed: %v", err)
		return response.ErrResp(err, codeCollectFailed)
	}

	return nil
}

// UncollectBlog 取消收藏帖子
func (l *BlogLogic) UncollectBlog(ctx context.Context, req types.UncollectBlogReq) error {
	defer utils.RecordTime(time.Now())()

	err := repo.NewBlogRepo(global.DB).UncollectBlog(int64(req.UserID), req.BlogID)
	if err != nil {
		zlog.CtxErrorf(ctx, "UncollectBlog failed:%v", err)
		return response.ErrResp(err, codeUncollectFailed)
	}

	return nil
}

// GetCollectedBlogs 获取用户收藏的帖子
func (l *BlogLogic) GetCollectedBlogs(ctx context.Context, req types.GetCollectedBlogsReq) (resp *types.GetCollectedBlogsResp, err error) {
	defer utils.RecordTime(time.Now())()

	blogs, err := repo.NewBlogRepo(global.DB).GetCollectedBlogs(int64(req.UserID))
	if err != nil {
		zlog.CtxErrorf(ctx, "GetCollectedBlogs failed: %v", err)
		return nil, response.ErrResp(err, codeGetCollectedFailed)
	}

	resp = &types.GetCollectedBlogsResp{
		List: blogs,
	}
	return resp, nil
}

// LikeBlog 点赞帖子
func (l *BlogLogic) LikeBlog(ctx context.Context, req types.LikeBlogReq) (resp *types.LikeBlogResp, err error) {
	defer utils.RecordTime(time.Now())()

	err = repo.NewBlogRepo(global.DB).LikeBlog(int64(req.UserID), req.BlogID)
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
func (l *BlogLogic) UnlikeBlog(ctx context.Context, req types.UnlikeBlogReq) (resp *types.UnlikeBlogResp, err error) {
	defer utils.RecordTime(time.Now())()

	err = repo.NewBlogRepo(global.DB).UnlikeBlog(int64(req.UserID), req.BlogID)
	if err != nil {
		zlog.CtxErrorf(ctx, "UnlikeBlog failed: %v", err)
		return nil, response.ErrResp(err, codeUnlikeFailed)
	}

	resp = &types.UnlikeBlogResp{
		BlogID:  req.BlogID,
		Success: true,
	}
	return resp, nil
}
