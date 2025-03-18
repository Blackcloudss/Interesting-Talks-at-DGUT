package types

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"mime/multipart"
)

// CreateBlogReq 创建帖子请求体
type CreateBlogReq struct {
	Content        string                  `json:"content"`
	Tag            string                  `json:"tag"`             // 帖子分区（主标签）
	SubTag         string                  `json:"sub_tag"`         // 子标签
	ViewPermission string                  `json:"view_permission"` // 访问权限，默认为“所有人”
	ImageFiles     []*multipart.FileHeader `form:"image_files"`     // 图片文件
}

// CreateBlogResp 创建帖子响应体
type CreateBlogResp struct {
	Blog   model.Blog    `json:"blog"`
	Images []model.Image `json:"images"`
}

// UpdateBlogReq 更新帖子请求体
type UpdateBlogReq struct {
	ID             int64                   `json:"id"`
	Content        string                  `json:"content"`
	Tag            string                  `json:"tag"`
	SubTag         string                  `json:"sub_tag"`
	ViewPermission string                  `json:"view_permission"`
	ImageFiles     []*multipart.FileHeader `form:"image_files"` // 图片文件
}

// UpdateBlogResp 更新帖子响应体
type UpdateBlogResp struct {
	Blog   model.Blog    `json:"blog"`
	Images []model.Image `json:"images"`
}

// DeleteBlogReq 删除帖子请求体
type DeleteBlogReq struct {
	ID int64 `json:"id"`
}

// DeleteBlogResp 删除帖子响应体
type DeleteBlogResp struct {
	Success bool `json:"success"`
}

// GetBlogByIDReq 获取帖子详情请求体
type GetBlogByIDReq struct {
	ID int64 `json:"id"`
}

// GetBlogByIDResp 获取帖子详情响应体
type GetBlogByIDResp struct {
	Blog   model.Blog     `json:"blog"`
	Images []*model.Image `json:"images"` // 确保返回的是图片列表
}

// GetBlogsReq 分页显示帖子请求体
type GetBlogsReq struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// GetBlogsResp 分页显示帖子响应体
type GetBlogsResp struct {
	Blogs    []model.Blog  `json:"list"`
	Images   []model.Image `json:"images"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

// GetBlogsByTagReq 根据标签显示帖子列表请求体
type GetBlogsByTagReq struct {
	SubTag string `json:"subtag"`
}

// GetBlogsByTagResp 根据标签显示帖子列表响应体
type GetBlogsByTagResp struct {
	Blogs  []model.Blog  `json:"list"`
	Images []model.Image `json:"images"`
}

// GetMyBlogsReq 获取当前用户发布的帖子请求体
type GetMyBlogsReq struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// GetMyBlogsResp 获取当前用户发布的帖子响应体
type GetMyBlogsResp struct {
	Blogs    []model.Blog  `json:"list"`
	Images   []model.Image `json:"images"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

// 获取其他用户帖子
type GetBlogsByUserIDReq struct {
	UserID   int64 `json:"user_id"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

// GetMyBlogsResp 获取其他用户发布的帖子响应体
type GetBlogsByUserIDResp struct {
	Blogs    []model.Blog  `json:"list"`
	Images   []model.Image `json:"images"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

// CollectBlogReq 收藏帖子请求体
type CollectBlogReq struct {
	BlogID int64 `json:"blog_id"`
}

// CollectBlogResp 收藏帖子响应体
type CollectBlogResp struct {
	Collection model.Collection `json:"collection"`
	Success    bool             `json:"success"`
}

// UncollectBlogReq 取消收藏帖子请求体
type UncollectBlogReq struct {
	BlogID int64 `json:"blog_id"`
}

// UncollectBlogResp 取消收藏帖子响应体
type UncollectBlogResp struct {
	Success bool `json:"success"`
}

// GetCollectedBlogsReq 获取用户收藏的帖子请求体
type GetCollectedBlogsReq struct {
	//UserID int `json:"user_id"`
}

// GetCollectedBlogsResp 获取用户收藏的帖子响应体
type GetCollectedBlogsResp struct {
	Blogs  []model.Blog  `json:"list"`
	Images []model.Image `json:"images"`
}

// LikeBlogReq 点赞帖子请求体
type LikeBlogReq struct {
	BlogID int64 `json:"blog_id"`
}

// LikeBlogResp 点赞帖子响应体
type LikeBlogResp struct {
	Success bool       `json:"success"`
	Like    model.Like `json:"like"`
}

// UnlikeBlogReq 取消点赞请求体
type UnlikeBlogReq struct {
	BlogID int64 `json:"blog_id"`
}

// UnlikeBlogResp 取消点赞响应体
type UnlikeBlogResp struct {
	Success bool `json:"success"`
}
