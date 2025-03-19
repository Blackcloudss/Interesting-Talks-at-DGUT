package api

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/jwt"
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
	UserID := jwt.GetUserId(c)
	resp, err := logic.NewCommentLogic().CreateComment(ctx, req, UserID)
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
	UserID := jwt.GetUserId(c)
	resp, err := logic.NewCommentLogic().DeleteComment(ctx, req, UserID)
	response.Response(c, resp, err)
	return
}

// 获取评论列表
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

// LikeCommentHandler 点赞评论
func LikeCommentHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.LikeCommentReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "LikeComment request error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "LikeComment request: %+v", req)
	UserID := jwt.GetUserId(c)
	resp, err := logic.NewCommentLogic().LikeComment(ctx, req, UserID)
	response.Response(c, resp, err)
	return
}

// UnlikeCommentHandler 取消点赞评论
func UnlikeCommentHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.UnlikeCommentReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "UnlikeComment request error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "UnlikeComment request: %+v", req)
	UserID := jwt.GetUserId(c)
	resp, err := logic.NewCommentLogic().UnlikeComment(ctx, req, UserID)
	response.Response(c, resp, err)
	return
}

//给评论回复
