package types

import "github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"

// SearchBlogsReq 搜索帖子请求体
type SearchBlogsReq struct {
	PageReq
	Keyword    string `json:"keyword" form:"keyword"`         // 搜索关键词
	SearchType string `json:"search_type" form:"search_type"` // 搜索类型: "title", "tag", "nickname", "all" (默认)
}

// SearchBlogsResp 搜索帖子响应体
type SearchBlogsResp struct {
	Blogs []BlogResp `json:"blogs"` // 帖子列表
	Total int64      `json:"total"` // 总帖子数
	PageReq
}

// GetSearchHistoryReq 获取搜索历史请求体
type GetSearchHistoryReq struct {
	PageReq
}

// GetSearchHistoryResp 获取搜索历史响应体
type GetSearchHistoryResp struct {
	History []model.SearchHistory `json:"history"` // 搜索历史列表
	Total   int64                 `json:"total"`   // 总记录数
	PageReq
}

type DeleteSearchHistoryReq struct {
	HistoryID int64 `json:"history_id" uri:"history_id"` // 搜索历史记录ID
}

type DeleteSearchHistoryResp struct{}

// 获取热搜榜
type GetHotSearchReq struct{}

type GetHotSearchResp struct {
	HotSearchList []string `json:"hot_search_list"` // 热搜榜列表
}
