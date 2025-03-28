package api

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/gin-gonic/gin"
)

// CreateNoticeHandler 创建公告
func CreateNoticeHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.CreateNoticeReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "CreateNotice request error: %v", err)
		return
	}
	zlog.CtxInfof(ctx, "CreateNotice request: %+v", req)
	resp, err := logic.NewNoticeLogic().CreateNotice(ctx, req)
	response.Response(c, resp, err)
	return
}

// UpdateNoticeHandler 更新公告
func UpdateNoticeHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.UpdateNoticeReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "UpdateNotice request error: %v", err)
		return
	}
	zlog.CtxInfof(ctx, "UpdateNotice request: %+v", req)
	resp, err := logic.NewNoticeLogic().UpdateNotice(ctx, req)
	response.Response(c, resp, err)
	return
}

// DeleteNoticeHandler 删除公告
func DeleteNoticeHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.DeleteNoticeReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "DeleteNotice request error: %v", err)
		return
	}
	zlog.CtxInfof(ctx, "DeleteNotice request: %+v", req)
	resp, err := logic.NewNoticeLogic().DeleteNotice(ctx, req)
	response.Response(c, resp, err)
	return
}

// GetNoticeHandler 获取公告详情
func GetNoticeHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.GetNoticeReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetNotice request error: %v", err)
		return
	}
	zlog.CtxInfof(ctx, "GetNotice request: %+v", req)
	resp, err := logic.NewNoticeLogic().GetNoticeByID(ctx, req)
	response.Response(c, resp, err)
	return
}
