package api

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/gin-gonic/gin"
)

// CreateCommentHandler 创建评论
func CreateComment(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.CreateCommentReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "CreateComment request error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "CreateComment request: %+v", req)
	resp, err := logic.NewCommentLogic().CreateComment(ctx, req)
	response.Response(c, resp, err)
	return
}

// DeleteCommentHandler 删除评论
func DeleteComment(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.DeleteCommentReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "DeleteComment request error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "DeleteComment request: %+v", req)
	resp, err := logic.NewCommentLogic().DeleteComment(ctx, req)
	response.Response(c, resp, err)
	return
}

func GetCommentList(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetCommentListReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetCommentList request error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "GetCommentList request: %+v", req)
	resp, err := logic.NewCommentLogic().GetCommentList(ctx, req)
	response.Response(c, resp, err)
	return
}
