package api

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/gin-gonic/gin"
)

// CreateCommentHandler 创建评论
func CreateCommentHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req := new(types.CreateCommentReq)
	if err := c.ShouldBindJSON(req); err != nil {
		zlog.CtxErrorf(ctx, "CreateComment request error: %v", err)
		response.Response(c, nil, err)
		return
	}

	zlog.CtxInfof(ctx, "CreateComment request: %+v", req)
	resp, err := logic.NewCommentLogic().CreateComment(ctx, *req)
	response.Response(c, resp, err)
}

// DeleteCommentHandler 删除评论
func DeleteCommentHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req := new(types.DeleteCommentReq)
	if err := c.ShouldBindJSON(req); err != nil {
		zlog.CtxErrorf(ctx, "DeleteComment request error: %v", err)
		response.Response(c, nil, err)
		return
	}

	zlog.CtxInfof(ctx, "DeleteComment request: %+v", req)
	resp, err := logic.NewCommentLogic().DeleteComment(ctx, req)
	response.Response(c, resp, err)
}

// GetCommentListHandler 获取评论列表
func GetCommentListHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req := new(types.GetCommentListReq)
	if err := c.ShouldBindJSON(req); err != nil {
		zlog.CtxErrorf(ctx, "GetCommentList request error: %v", err)
		response.Response(c, nil, err)
		return
	}

	zlog.CtxInfof(ctx, "GetCommentList request: %+v", req)
	resp, err := logic.NewCommentLogic().GetCommentList(ctx, *req)
	response.Response(c, resp, err)
}
