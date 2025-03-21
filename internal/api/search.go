package api

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/jwt"
	"github.com/gin-gonic/gin"
)

// 模糊搜索
// SearchBlogsHandler 搜索帖子
func SearchBlogsHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.SearchBlogsReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "SearchBlogs request error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}

	UserID := jwt.GetUserId(c)

	resp, err := logic.NewSearchLogic().SearchBlogs(ctx, req, UserID)
	response.Response(c, resp, err)
}

// GetSearchHistoryHandler 获取用户的搜索历史
func GetSearchHistoryHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetSearchHistoryReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetSearchHistory request error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}

	UserID := jwt.GetUserId(c)

	resp, err := logic.NewSearchLogic().GetSearchHistory(ctx, req, UserID)
	response.Response(c, resp, err)
}

//删除搜索历史
