package controller

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/gin-gonic/gin"
	"strconv"
)

// CreateBlogHandler 创建帖子
func CreateBlogHandler(c *gin.Context) {
	// 从 Gin 中获取上下文
	ctx := zlog.GetCtxFromGin(c)

	// 1. 获取参数及校验参数
	var blog model.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		zlog.CtxErrorf(ctx, "c.ShouldBindJSON(blog) err: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID) // 参数无效
		return
	}

	// 2. 获取作者ID（当前请求的UserID）
	userID, err := getCurrentUserID(c)
	if err != nil {
		zlog.CtxErrorf(ctx, "getCurrentUserID() failed: %v", err)
		response.NewResponse(c).Error(response.USER_NOT_LOGIN) // 用户未登录
		return
	}
	blog.UserID = userID

	// 3. 使用雪花算法生成帖子ID
	blogID := global.Node.Generate().Int64()
	blog.ID = uint64(blogID)

	// 4. 创建帖子
	if err := logic.NewBlogLogic().CreateBlog(ctx, &blog); err != nil {
		zlog.CtxErrorf(ctx, "logic.CreateBlog failed: %v", err)
		response.NewResponse(c).Error(response.INTERNAL_ERROR) // 创建失败
		return
	}

	// 5. 返回响应
	zlog.CtxInfof(ctx, "Blog created successfully: %+v", blog)
	response.NewResponse(c).Success(blog) // 创建成功
}

// UpdateBlogHandler 更新帖子
func UpdateBlogHandler(c *gin.Context) {
	// 从 Gin 中获取上下文
	ctx := zlog.GetCtxFromGin(c)

	// 1. 获取帖子ID并校验
	id := c.Param("id")
	blogID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "strconv.ParseUint(id) err: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID) // 参数无效
		return
	}

	// 2. 获取更新数据并校验
	var blog model.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		zlog.CtxErrorf(ctx, "c.ShouldBindJSON(blog) err: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID) // 参数无效
		return
	}
	blog.ID = blogID

	// 3. 检查权限
	userID, err := getCurrentUserID(c)
	if err != nil {
		zlog.CtxErrorf(ctx, "getCurrentUserID() failed: %v", err)
		response.NewResponse(c).Error(response.USER_NOT_LOGIN) // 用户未登录
		return
	}
	if blog.UserID != userID {
		zlog.CtxErrorf(ctx, "User has no permission to update this blog")
		response.NewResponse(c).Error(response.INSUFFICENT_PERMISSIONS) // 权限不足
		return
	}

	// 4. 更新帖子
	if err := logic.NewBlogLogic().UpdateBlog(ctx, &blog); err != nil {
		zlog.CtxErrorf(ctx, "logic.UpdateBlog failed: %v", err)
		response.NewResponse(c).Error(response.INTERNAL_ERROR) // 更新失败
		return
	}

	// 5. 返回响应
	zlog.CtxInfof(ctx, "Blog updated successfully: %+v", blog)
	response.NewResponse(c).Success(blog) // 更新成功
}

// DeleteBlogHandler 删除帖子
func DeleteBlogHandler(c *gin.Context) {
	// 从 Gin 中获取上下文
	ctx := zlog.GetCtxFromGin(c)

	// 1. 获取帖子ID并校验
	id := c.Param("id")
	blogID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "strconv.ParseUint(id) err: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID) // 参数无效
		return
	}

	// 2. 获取当前用户id
	userID, err := getCurrentUserID(c)
	if err != nil {
		zlog.CtxErrorf(ctx, "getCurrentUserID() failed: %v", err)
		response.NewResponse(c).Error(response.USER_NOT_LOGIN) // 用户未登录
		return
	}

	// 3. 获取帖子详情
	blog, err := logic.NewBlogLogic().GetBlogByID(ctx, blogID, userID)
	if err != nil {
		zlog.CtxErrorf(ctx, "logic.GetBlogByID failed: %v", err)
		response.NewResponse(c).Error(response.MESSAGE_NOT_EXIST) // 帖子不存在
		return
	}

	// 检查权限
	if blog.UserID != userID {
		zlog.CtxErrorf(ctx, "User has no permission to delete this blog")
		response.NewResponse(c).Error(response.INSUFFICENT_PERMISSIONS) // 权限不足
		return
	}

	// 4. 删除帖子
	if err := logic.NewBlogLogic().DeleteBlog(ctx, blogID); err != nil {
		zlog.CtxErrorf(ctx, "logic.DeleteBlog failed: %v", err)
		response.NewResponse(c).Error(response.INTERNAL_ERROR) // 删除失败
		return
	}

	// 5. 返回响应
	zlog.CtxInfof(ctx, "Blog deleted successfully: %+v", blog)
	response.NewResponse(c).Success(nil) // 删除成功
}

// GetBlogByIDHandler 获取帖子详情
func GetBlogByIDHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)

	blogID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "Invalid blog ID: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID) // 参数无效
		return
	}

	userID, err := getCurrentUserID(c)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get current user ID: %v", err)
		response.NewResponse(c).Error(response.USER_NOT_LOGIN) // 用户未登录
		return
	}

	blog, err := logic.NewBlogLogic().GetBlogByID(ctx, blogID, userID)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get blog (blogID: %d): %v", blogID, err)
		response.NewResponse(c).Error(response.MESSAGE_NOT_EXIST) // 帖子不存在或无权限
		return
	}

	zlog.CtxInfof(ctx, "Blog retrieved successfully: %+v", blog)
	response.NewResponse(c).Success(blog) // 获取成功
}

// GetBlogsHandler 分页显示帖子（支持下拉刷新）
func GetBlogsHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)

	// 获取分页参数
	page, pageSize := paginate(c)

	// 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get current user ID: %v", err)
		response.NewResponse(c).Error(response.USER_NOT_LOGIN) // 用户未登录
		return
	}

	// 获取最新帖子的ID（用于下拉刷新）
	latestIDStr := c.Query("latestID")
	var latestID uint64
	if latestIDStr != "" {
		latestID, err = strconv.ParseUint(latestIDStr, 10, 64)
		if err != nil {
			zlog.CtxErrorf(ctx, "Invalid latestID: %v", err)
			response.NewResponse(c).Error(response.PARAM_NOT_VALID) // 参数无效
			return
		}
	}

	var blogs []model.Blog
	var total int64

	// 根据 latestID 判断是普通分页还是下拉刷新
	if latestID > 0 {
		blogs, err = logic.NewBlogLogic().GetBlogsAfterID(ctx, latestID, pageSize)
		if err != nil {
			zlog.CtxErrorf(ctx, "Failed to get blogs after ID: %v", err)
			response.NewResponse(c).Error(response.INTERNAL_ERROR) // 获取失败
			return
		}
		total = int64(len(blogs)) // 下拉刷新时，total 为实际获取的数量
	} else {
		blogs, total, err = logic.NewBlogLogic().GetBlogs(ctx, page, pageSize, userID)
		if err != nil {
			zlog.CtxErrorf(ctx, "Failed to get blogs: %v", err)
			response.NewResponse(c).Error(response.INTERNAL_ERROR) // 获取失败
			return
		}
	}

	// 返回成功响应
	zlog.CtxInfof(ctx, "Blogs retrieved successfully: %+v", blogs)
	response.NewResponse(c).Success(gin.H{
		"list":     blogs,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetBlogsByTagHandler 根据标签显示帖子列表
func GetBlogsByTagHandler(c *gin.Context) {
	// 从 Gin 中获取上下文
	ctx := zlog.GetCtxFromGin(c)

	// 1. 获取标签和分页参数
	tag := c.Query("tag")
	page, pageSize := paginate(c)

	// 2. 获取帖子列表
	blogs, total, err := logic.NewBlogLogic().GetBlogsByTag(ctx, tag, page, pageSize)
	if err != nil {
		zlog.CtxErrorf(ctx, "logic.GetBlogsByTag failed: %v", err)
		response.NewResponse(c).Error(response.INTERNAL_ERROR) // 获取失败
		return
	}

	// 3. 返回响应
	zlog.CtxInfof(ctx, "Blogs retrieved successfully: %+v", blogs)
	response.NewResponse(c).Success(gin.H{
		"list":     blogs,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}) // 获取成功
}

// GetMyBlogsHandler 获取当前用户发布的帖子
func GetMyBlogsHandler(c *gin.Context) {
	// 从 Gin 中获取上下文
	ctx := zlog.GetCtxFromGin(c)

	// 1. 获取分页参数
	page, pageSize := paginate(c)

	// 2. 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		zlog.CtxErrorf(ctx, "getCurrentUserID() failed: %v", err)
		response.NewResponse(c).Error(response.USER_NOT_LOGIN) // 用户未登录
		return
	}

	// 3. 获取帖子列表
	blogs, total, err := logic.NewBlogLogic().GetBlogsByUserID(ctx, userID, page, pageSize)
	if err != nil {
		zlog.CtxErrorf(ctx, "logic.GetBlogsByUserID failed: %v", err)
		response.NewResponse(c).Error(response.INTERNAL_ERROR) // 获取失败
		return
	}

	// 4. 返回响应
	zlog.CtxInfof(ctx, "Blogs retrieved successfully: %+v", blogs)
	response.NewResponse(c).Success(gin.H{
		"list":     blogs,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}) // 获取成功
}
