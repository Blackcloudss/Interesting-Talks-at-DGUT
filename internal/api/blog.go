package api

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/image"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/jwt"
	"github.com/gin-gonic/gin"
)

// CreateBlogHandler 创建帖子
func CreateBlogHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.CreateBlogReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "CreateBlog request error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}

	// 处理图片上传（可能有多张）
	var imageUrls []string
	if req.ImageFiles != nil {
		for _, fileHeader := range req.ImageFiles {
			imageUrl, err := image.UploadImage(fileHeader)
			if err != nil {
				zlog.CtxErrorf(ctx, "Upload image error: %v", err)
				response.NewResponse(c).Error(response.INTERNAL_ERROR)
				return
			}
			zlog.CtxInfof(ctx, "Image uploaded successfully, URL: %s", imageUrl)
			imageUrls = append(imageUrls, imageUrl)
		}
	}
	UserID := jwt.GetUserId(c)
	zlog.CtxInfof(ctx, "CreateBlog request: %+v", req)
	resp, err := logic.NewBlogLogic().CreateBlog(ctx, req, UserID, imageUrls)
	response.Response(c, resp, err)
	return
}

// UpdateBlogHandler 更新帖子
func UpdateBlogHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.UpdateBlogReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "UpdateBlog request error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	// 处理图片上传（可能有多张）
	var imageUrls []string
	if req.ImageFiles != nil {
		for _, fileHeader := range req.ImageFiles {
			imageUrl, err := image.UploadImage(fileHeader)
			if err != nil {
				zlog.CtxErrorf(ctx, "Upload image error: %v", err)
				response.NewResponse(c).Error(response.INTERNAL_ERROR)
				return
			}
			zlog.CtxInfof(ctx, "Image uploaded successfully, URL: %s", imageUrl)
			imageUrls = append(imageUrls, imageUrl)
		}
	}

	zlog.CtxInfof(ctx, "UpdateBlog request: %+v", req)
	resp, err := logic.NewBlogLogic().UpdateBlog(ctx, req, imageUrls)
	response.Response(c, resp, err)
	return
}

// DeleteBlogHandler 删除帖子
func DeleteBlogHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.DeleteBlogReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "DeleteBlog request error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "DeleteBlog request: %+v", req)
	resp, err := logic.NewBlogLogic().DeleteBlog(ctx, req)
	response.Response(c, resp, err)
	return
}

// GetBlogByIDHandler 获取帖子详情
func GetBlogByIDHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetBlogByIDReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetBlogByID request error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "GetBlogByID request: %+v", req)
	resp, err := logic.NewBlogLogic().GetBlogByID(ctx, req)
	response.Response(c, resp, err)
	return
}

// GetBlogsHandler 分页显示帖子
func GetBlogsHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetBlogsReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetBlogs request error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "GetBlogs request: %+v", req)
	resp, err := logic.NewBlogLogic().GetBlogs(ctx, req)
	response.Response(c, resp, err)
	return
}

// GetBlogsByTagHandler 根据标签显示帖子列表
func GetBlogsByTagHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetBlogsByTagReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetBlogsByTag request error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "GetBlogsByTag request: %+v", req)
	resp, err := logic.NewBlogLogic().GetBlogsByTag(ctx, req)
	response.Response(c, resp, err)
	return
}

// GetMyBlogsHandler 获取当前用户发布的帖子
func GetMyBlogsHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetMyBlogsReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetMyBlogs request error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "GetMyBlogs request: %+v", req)

	UserID := jwt.GetUserId(c)
	resp, err := logic.NewBlogLogic().GetMyBlogs(ctx, req, UserID)
	response.Response(c, resp, err)
	return
}

// GetMyBlogsHandler 获取其他用户发布的帖子
func GetBlogsByUserIDHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetBlogsByUserIDReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetBlogsByUserID request error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "GetBlogsByUserID request: %+v", req)
	resp, err := logic.NewBlogLogic().GetBlogsByUserID(ctx, req)
	response.Response(c, resp, err)
	return
}

// CollectBlogHandler 收藏帖子
func CollectBlogHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.CollectBlogReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "CollectBlog request error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "CollectBlog request: %+v", req)

	UserID := jwt.GetUserId(c)
	resp, err := logic.NewBlogLogic().CollectBlog(ctx, req, UserID)
	response.Response(c, resp, err)
	return
}

// UncollectBlogHandler 取消收藏帖子
func UncollectBlogHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.UncollectBlogReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "UncollectBlog request error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "UncollectBlog request: %+v", req)
	UserID := jwt.GetUserId(c)
	resp, err := logic.NewBlogLogic().UncollectBlog(ctx, req, UserID)
	response.Response(c, resp, err)
	return
}

// GetCollectedBlogsHandler 获取用户收藏的帖子
func GetCollectedBlogsHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetCollectedBlogsReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetCollectedBlogs request error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "GetCollectedBlogs request: %+v", req)
	UserID := jwt.GetUserId(c)
	resp, err := logic.NewBlogLogic().GetCollectedBlogs(ctx, req, UserID)
	response.Response(c, resp, err)
	return
}

// LikeBlogHandler 点赞帖子
func LikeBlogHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.LikeBlogReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "LikeBlog request error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "LikeBlog request: %+v", req)
	UserID := jwt.GetUserId(c)
	resp, err := logic.NewBlogLogic().LikeBlog(ctx, req, UserID)
	response.Response(c, resp, err)
	return
}

// UnlikeBlogHandler 取消点赞
func UnlikeBlogHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.UnlikeBlogReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "UnlikeBlog request error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "UnlikeBlog request: %+v", req)
	UserID := jwt.GetUserId(c)
	resp, err := logic.NewBlogLogic().UnlikeBlog(ctx, req, UserID)
	response.Response(c, resp, err)
	return
}
