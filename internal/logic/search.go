package logic

import (
	"context"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/repo"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils"
	"time"
)

type SearchLogic struct{}

func NewSearchLogic() *SearchLogic {
	return &SearchLogic{}
}

// SearchBlogs 搜索帖子
func (l *SearchLogic) SearchBlogs(ctx context.Context, req types.SearchBlogsReq, UserID int64) (resp *types.SearchBlogsResp, err error) {
	defer utils.RecordTime(time.Now())()

	// 搜索帖子
	blogs, total, err := repo.NewSearchRepo(global.DB).SearchBlogs(req.Keyword, req.SearchType, req.Page, req.PageSize)
	if err != nil {
		zlog.CtxErrorf(ctx, "SearchBlogs failed: %v", err)
		return nil, err
	}

	// 记录搜索历史
	if err := repo.NewSearchRepo(global.DB).RecordSearchHistory(UserID, req.Keyword); err != nil {
		zlog.CtxErrorf(ctx, "RecordSearchHistory failed: %v", err)
		return nil, err
	}

	resp = &types.SearchBlogsResp{
		Blogs: blogs,
		Total: total,
		PageReq: types.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}
	return resp, nil
}

// GetSearchHistory 获取用户的搜索历史
func (l *SearchLogic) GetSearchHistory(ctx context.Context, req types.GetSearchHistoryReq, UserID int64) (resp *types.GetSearchHistoryResp, err error) {
	defer utils.RecordTime(time.Now())()

	// 获取用户的搜索历史
	history, total, err := repo.NewSearchRepo(global.DB).GetSearchHistory(UserID, req.Page, req.PageSize)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetSearchHistory failed: %v", err)
		return nil, err
	}

	resp = &types.GetSearchHistoryResp{
		History: history,
		Total:   total,
		PageReq: types.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}
	return resp, nil
}

func (l *SearchLogic) DeleteSearchHistory(ctx context.Context, req types.DeleteSearchHistoryReq, UserID int64) (resp *types.DeleteSearchHistoryResp, err error) {
	defer utils.RecordTime(time.Now())()

	// 删除指定的搜索历史记录
	err = repo.NewSearchRepo(global.DB).DeleteSearchHistory(UserID, req.HistoryID)
	if err != nil {
		zlog.CtxErrorf(ctx, "DeleteSearchHistory failed: %v", err)
		return nil, err
	}

	return resp, nil
}
