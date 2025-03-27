package types

import (
	"mime/multipart"
	"time"
)

// CreateBlogReq 创建帖子请求体
type CreateBlogReq struct {
	Title          string                  `form:"title"`
	Content        string                  `form:"content"`
	BlogTag        string                  `form:"blog_tag"`
	SubTag         string                  `form:"sub_tag"`
	ViewPermission string                  `form:"view_permission"`
	ImageFiles     []*multipart.FileHeader `form:"image_files"`
}

// CreateBlogResp 创建帖子响应体
type CreateBlogResp struct {
	BlogID   int64     `json:"blog_id"`
	CreateAt time.Time `json:"create_at"`
}

// UpdateBlogReq 更新帖子请求体
type UpdateBlogReq struct {
	ID             int64                   `form:"id"`
	Title          string                  `form:"title"`
	Content        string                  `form:"content"`
	BlogTag        string                  `form:"blog_tag"`
	SubTag         string                  `form:"sub_tag"`
	ViewPermission string                  `form:"view_permission"`
	ImageFiles     []*multipart.FileHeader `form:"image_files"`
}

// UpdateBlogResp 更新帖子响应体
type UpdateBlogResp struct {
	UpdateAt time.Time `json:"update_at"`
}

// DeleteBlogReq 删除帖子请求体
type DeleteBlogReq struct {
	ID int64 `json:"blog_id"`
}

// DeleteBlogResp 删除帖子响应体
type DeleteBlogResp struct {
}

// 获取帖子内容通用响应体
type BlogResp struct {
	BlogID         int64     `json:"blog_id"`
	CreateAt       time.Time `json:"create_at"`
	UpdateAt       time.Time `json:"update_at"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	BeLiked        int       ` json:"be_liked"  `     // 点赞数，默认为0
	BeCollected    int       `json:"be_collected"`    // 收藏数，默认为0
	CommentCount   int       ` json:"comment_count"`  // 评论数，默认为0（字段名更清晰）
	BlogTag        string    ` json:"blog_tag"`       // 帖子分区（主标签）
	SubTag         string    ` json:"sub_tag"`        // 子标签
	ViewPermission string    `json:"view_permission"` // 访问权限，默认为“所有人”
	UserID         int64     `json:"user_id"`         // 用户ID
	Nickname       string    `json:"nickname"`        // 昵称
	Avatar         string    `json:"avatar"`          // 头像
	Tag            string    `json:"tag"`             // 标签
}

// GetBlogByIDReq 获取帖子详情请求体
type GetBlogByIDReq struct {
	ID int64 `json:"id" form:"id" `
}

// GetBlogByIDResp 获取帖子详情响应体
type GetBlogByIDResp struct {
	BlogResp
	Images []string `json:"images"` // 帖子的图片列表路径
}

// GetBlogsReq 分页显示帖子请求体
type GetBlogsReq struct {
	PageReq
}

// GetBlogsResp 分页显示帖子响应体
type GetBlogsResp struct {
	Blogs []BlogResp `json:"blogs"` // 帖子列表
	//Images map[int64][]string `json:"images"` // 每个帖子的图片列表
	Total int64 `json:"total"` // 总帖子数
	PageReq
}

// GetBlogsByTagReq 根据标签显示帖子列表请求体
type GetBlogsByTagReq struct {
	SubTag string `json:"sub_tag" form:"sub_tag"`
	PageReq
}

// GetBlogsByTagResp 根据标签显示帖子列表响应体
type GetBlogsByTagResp struct {
	Blogs []BlogResp `json:"blogs"` // 帖子列表
	//Images map[int64][]string `json:"images"` // 每个帖子的图片列表
	Total int64 `json:"total"` // 总帖子数
	PageReq
}

// GetMyBlogsReq 获取当前用户发布的帖子请求体
type GetMyBlogsReq struct {
	PageReq
}

// GetMyBlogsResp 获取当前用户发布的帖子响应体
type GetMyBlogsResp struct {
	Blogs []BlogResp `json:"blogs"` // 帖子列表
	//Images map[int64][]string `json:"images"` // 每个帖子的图片列表
	Total int64 `json:"total"` // 总帖子数
	PageReq
}

// GetBlogsByUserIDReq 获取其他用户发布的帖子请求体
type GetBlogsByUserIDReq struct {
	UserID int64 `json:"user_id" form:"user_id"`
	PageReq
}

// GetBlogsByUserIDResp 获取其他用户发布的帖子响应体
type GetBlogsByUserIDResp struct {
	Blogs []BlogResp `json:"blogs"` // 帖子列表
	//Images map[int64][]string `json:"images"` // 每个帖子的图片列表
	Total int64 `json:"total"` // 总帖子数
	PageReq
}

// CollectBlogReq 收藏帖子请求体
type CollectBlogReq struct {
	BlogID int64 `json:"blog_id"`
}

// CollectBlogResp 收藏帖子响应体
type CollectBlogResp struct {
}

// UncollectBlogReq 取消收藏帖子请求体
type UncollectBlogReq struct {
	BlogID int64 `json:"blog_id"`
}

// UncollectBlogResp 取消收藏帖子响应体
type UncollectBlogResp struct {
}

// GetCollectedBlogsReq 获取用户收藏的帖子请求体
type GetCollectedBlogsReq struct {
	PageReq
}

// GetCollectedBlogsResp 获取用户收藏的帖子响应体
type GetCollectedBlogsResp struct {
	Blogs []BlogResp `json:"blogs"` // 帖子列表
	//Images map[int64][]string `json:"images"` // 每个帖子的图片列表
	Total int64 `json:"total"` // 总帖子数
	PageReq
}

// LikeBlogReq 点赞帖子请求体
type LikeBlogReq struct {
	BlogID int64 `json:"blog_id"`
}

// LikeBlogResp 点赞帖子响应体
type LikeBlogResp struct {
}

// UnlikeBlogReq 取消点赞请求体
type UnlikeBlogReq struct {
	BlogID int64 `json:"blog_id"`
}

// UnlikeBlogResp 取消点赞响应体
type UnlikeBlogResp struct {
}
