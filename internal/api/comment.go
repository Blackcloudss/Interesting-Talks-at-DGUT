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

// LikeComment 点赞评论
func LikeComment(c *gin.Context) {
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

// UnlikeComment 取消点赞评论
func UnlikeComment(c *gin.Context) {
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

// 获取评论列表，按发布时间排序（分页）
// 显示一级评论时要显示部分二级评论，点展开更多显示更多二级评论

// 获取评论列表（按点赞数排序）（分页）
