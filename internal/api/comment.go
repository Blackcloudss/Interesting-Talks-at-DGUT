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
		return
	}
	zlog.CtxInfof(ctx, "DeleteComment request: %+v", req)
	resp, err := logic.NewCommentLogic().DeleteComment(ctx, req)
	response.Response(c, resp, err)
	return
}

// LikeComment 点赞评论
func LikeComment(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.LikeCommentReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "LikeComment request error: %v", err)
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
		return
	}
	zlog.CtxInfof(ctx, "UnlikeComment request: %+v", req)
	UserID := jwt.GetUserId(c)
	resp, err := logic.NewCommentLogic().UnlikeComment(ctx, req, UserID)
	response.Response(c, resp, err)
	return
}

// 获取评论列表，按发布时间排序或按点赞数排列（分页）
// 分页加载：通过滚动事件触发分页加载更多评论。
//展开二级评论：点击“展开更多”按钮，发送请求加载更多二级评论。

// GetCommentList 获取评论列表
func GetCommentList(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetCommentListReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetCommentList request error: %v", err)
		return
	}
	zlog.CtxInfof(ctx, "GetCommentList request: %+v", req)

	resp, err := logic.NewCommentLogic().GetCommentList(ctx, req)
	response.Response(c, resp, err)
}

// GetMoreSecondCommentHandler 获取更多二级评论
func GetMoreSecondComment(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetSecondCommentListReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetMoreSecondComment request error: %v", err)
		return
	}
	zlog.CtxInfof(ctx, "GetMoreSecondComment request: %+v", req)

	resp, err := logic.NewCommentLogic().GetSecondCommentList(ctx, req)
	response.Response(c, resp, err)
}
