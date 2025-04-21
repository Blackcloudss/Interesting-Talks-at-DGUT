package logic

import (
	"context"
	"errors"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/repo"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"gorm.io/gorm"
	"time"
)

// 定义内部逻辑错误
var (
	NOTICE_NOT_FOUND     = response.MsgCode{Code: 40041, Msg: "公告不存在"}
	NOTICE_CREATE_FAILED = response.MsgCode{Code: 40042, Msg: "创建公告失败"}
	NOTICE_UPDATE_FAILED = response.MsgCode{Code: 40043, Msg: "更新公告失败"}
	NOTICE_DELETE_FAILED = response.MsgCode{Code: 40044, Msg: "删除公告失败"}
	NOTICE_GET_FAILED    = response.MsgCode{Code: 40045, Msg: "获取公告失败"}
)

type NoticeLogic struct{}

// NewNoticeLogic 创建通知逻辑层实例
func NewNoticeLogic() *NoticeLogic {
	return &NoticeLogic{}
}

// CreateNotice 创建公告
func (l *NoticeLogic) CreateNotice(ctx context.Context, req types.CreateNoticeReq) (*types.CreateNoticeResp, error) {
	notice, err := repo.NewNoticeRepo(global.DB).CreateNotice(req)
	if err != nil {
		zlog.CtxErrorf(ctx, "CreateNotice failed: %v", err)
		return nil, response.ErrResp(err, NOTICE_CREATE_FAILED)
	}
	zlog.CtxInfof(ctx, "Notice created successfully (noticeID: %d)", notice.ID)
	return &types.CreateNoticeResp{
		ID:       notice.ID,
		CreateAt: notice.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// UpdateNotice 更新公告
func (l *NoticeLogic) UpdateNotice(ctx context.Context, req types.UpdateNoticeReq) (*types.UpdateNoticeResp, error) {
	if err := repo.NewNoticeRepo(global.DB).UpdateNotice(req); err != nil {
		zlog.CtxErrorf(ctx, "UpdateNotice failed: %v", err)
		return nil, response.ErrResp(err, NOTICE_UPDATE_FAILED)
	}
	zlog.CtxInfof(ctx, "Notice updated successfully (noticeID: %d)", req.ID)
	return &types.UpdateNoticeResp{
		UpdateAt: time.Now().Unix(),
	}, nil
}

// DeleteNotice 删除公告
func (l *NoticeLogic) DeleteNotice(ctx context.Context, req types.DeleteNoticeReq) (*types.DeleteNoticeResp, error) {
	if err := repo.NewNoticeRepo(global.DB).DeleteNotice(req.ID); err != nil {
		zlog.CtxErrorf(ctx, "DeleteNotice failed: %v", err)
		return nil, response.ErrResp(err, NOTICE_DELETE_FAILED)
	}
	zlog.CtxInfof(ctx, "Notice deleted successfully (noticeID: %d)", req.ID)
	return &types.DeleteNoticeResp{}, nil
}

// GetNoticeByID 获取公告详情
func (l *NoticeLogic) GetNoticeByID(ctx context.Context, req types.GetNoticeReq) (*types.GetNoticeResp, error) {
	notice, err := repo.NewNoticeRepo(global.DB).GetNoticeByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "GetNoticeByID failed: notice not found (noticeID: %d)", req.ID)
			return nil, response.ErrResp(err, NOTICE_NOT_FOUND)
		}
		zlog.CtxErrorf(ctx, "GetNoticeByID failed: %v", err)
		return nil, response.ErrResp(err, NOTICE_GET_FAILED)
	}
	zlog.CtxInfof(ctx, "Notice retrieved successfully (noticeID: %d)", req.ID)
	return &types.GetNoticeResp{
		UpdateAt: notice.UpdatedAt.Unix(),
		Content:  notice.Content,
	}, nil
}
