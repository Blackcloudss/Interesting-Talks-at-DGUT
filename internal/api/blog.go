package api

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/gin-gonic/gin"
)

// CreateBlogHandler 创建帖子
func CreateBlogHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.CreateBlogReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "CreateBlog request error: %v", err)
		response.Response(c, nil, err)
		return
	}
	zlog.CtxInfof(ctx, "CreateBlog request: %v", req)
	resp, err := logic.NewBlogLogic().CreateBlog(ctx, req)
	response.Response(c, resp, err)
}

// UpdateBlogHandler 更新帖子
func UpdateBlogHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.UpdateBlogReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "UpdateBlog request error: %v", err)
		response.Response(c, nil, err)
		return
	}
	zlog.CtxInfof(ctx, "UpdateBlog request: %v", req)
	resp, err := logic.NewBlogLogic().UpdateBlog(ctx, req)
	response.Response(c, resp, err)
}

// DeleteBlogHandler 删除帖子
func DeleteBlogHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.DeleteBlogReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "DeleteBlog request error: %v", err)
		response.Response(c, nil, err)
		return
	}
	zlog.CtxInfof(ctx, "DeleteBlog request: %v", req)
	err = logic.NewBlogLogic().DeleteBlog(ctx, req)
	response.Response(c, nil, err)
}

// GetBlogByIDHandler 获取帖子详情
func GetBlogByIDHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetBlogByIDReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetBlogByID request error: %v", err)
		response.Response(c, nil, err)
		return
	}
	zlog.CtxInfof(ctx, "GetBlogByID request: %v", req)
	resp, err := logic.NewBlogLogic().GetBlogByID(ctx, req)
	response.Response(c, resp, err)
}

// GetBlogsHandler 分页显示帖子
func GetBlogsHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetBlogsReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetBlogs request error: %v", err)
		response.Response(c, nil, err)
		return
	}
	zlog.CtxInfof(ctx, "GetBlogs request: %v", req)
	resp, err := logic.NewBlogLogic().GetBlogs(ctx, req)
	response.Response(c, resp, err)
}

// GetBlogsByTagHandler 根据标签显示帖子列表
func GetBlogsByTagHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetBlogsByTagReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetBlogsByTag request error: %v", err)
		response.Response(c, nil, err)
		return
	}
	zlog.CtxInfof(ctx, "GetBlogsByTag request: %v", req)
	resp, err := logic.NewBlogLogic().GetBlogsByTag(ctx, req)
	response.Response(c, resp, err)
}

// GetMyBlogsHandler 获取当前用户发布的帖子
func GetMyBlogsHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetMyBlogsReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetMyBlogs request error: %v", err)
		response.Response(c, nil, err)
		return
	}
	zlog.CtxInfof(ctx, "GetMyBlogs request: %v", req)
	resp, err := logic.NewBlogLogic().GetMyBlogs(ctx, req)
	response.Response(c, resp, err)
}

// CollectBlogHandler 收藏帖子
func CollectBlogHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.CollectBlogReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "CollectBlog request error: %v", err)
		response.Response(c, nil, err)
		return
	}
	zlog.CtxInfof(ctx, "CollectBlog request: %v", req)
	err = logic.NewBlogLogic().CollectBlog(ctx, req)
	response.Response(c, nil, err)
}

// UncollectBlogHandler 取消收藏帖子
func UncollectBlogHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.UncollectBlogReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "UncollectBlog request error: %v", err)
		response.Response(c, nil, err)
		return
	}
	zlog.CtxInfof(ctx, "UncollectBlog request: %v", req)
	err = logic.NewBlogLogic().UncollectBlog(ctx, req)
	response.Response(c, nil, err)
}

// GetCollectedBlogsHandler 获取用户收藏的帖子
func GetCollectedBlogsHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetCollectedBlogsReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetCollectedBlogs request error: %v", err)
		response.Response(c, nil, err)
		return
	}
	zlog.CtxInfof(ctx, "GetCollectedBlogs request: %v", req)
	resp, err := logic.NewBlogLogic().GetCollectedBlogs(ctx, req)
	response.Response(c, resp, err)
}

// LikeBlogHandler 点赞帖子
func LikeBlogHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.LikeBlogReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "LikeBlog request error: %v", err)
		response.Response(c, nil, err)
		return
	}
	zlog.CtxInfof(ctx, "LikeBlog request: %v", req)
	err = logic.NewBlogLogic().LikeBlog(ctx, req)
	response.Response(c, nil, err)
}

// UnlikeBlogHandler 取消点赞
func UnlikeBlogHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.UnlikeBlogReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "UnlikeBlog request error: %v", err)
		response.Response(c, nil, err)
		return
	}
	zlog.CtxInfof(ctx, "UnlikeBlog request: %v", req)
	err = logic.NewBlogLogic().UnlikeBlog(ctx, req)
	response.Response(c, nil, err)
}
