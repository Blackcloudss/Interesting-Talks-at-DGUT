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
	codeNoticeNotFound     = response.MsgCode{Code: 40041, Msg: "通知不存在"}
	codeNoticeCreateFailed = response.MsgCode{Code: 40042, Msg: "创建通知失败"}
	codeNoticeUpdateFailed = response.MsgCode{Code: 40043, Msg: "更新通知失败"}
	codeNoticeDeleteFailed = response.MsgCode{Code: 40044, Msg: "删除通知失败"}
)

type NoticeLogic struct{}

// NewNoticeLogic 创建通知逻辑层实例
func NewNoticeLogic() *NoticeLogic {
	return &NoticeLogic{}
}

// CreateNotice 创建公告
func (l *NoticeLogic) CreateNotice(ctx context.Context, req types.CreateNoticeReq, UserID int64) (*types.CreateNoticeResp, error) {
	noticeRepo := repo.NewNoticeRepo(global.DB)
	notice, err := noticeRepo.CreateNotice(req, UserID)
	if err != nil {
		zlog.CtxErrorf(ctx, "CreateNotice failed: %v", err)
		return nil, response.ErrResp(err, codeNoticeCreateFailed)
	}
	zlog.CtxInfof(ctx, "Notice created successfully (noticeID: %d)", notice.ID)
	return &types.CreateNoticeResp{
		ID:       notice.ID,
		CreateAt: notice.CreatedAt.Unix(),
	}, nil
}

// UpdateNotice 更新公告
func (l *NoticeLogic) UpdateNotice(ctx context.Context, req types.UpdateNoticeReq) (*types.UpdateNoticeResp, error) {
	noticeRepo := repo.NewNoticeRepo(global.DB)
	if err := noticeRepo.UpdateNotice(req); err != nil {
		zlog.CtxErrorf(ctx, "UpdateNotice failed: %v", err)
		return nil, response.ErrResp(err, codeNoticeUpdateFailed)
	}
	zlog.CtxInfof(ctx, "Notice updated successfully (noticeID: %d)", req.ID)
	return &types.UpdateNoticeResp{
		UpdateAt: time.Now().Unix(),
	}, nil
}

// DeleteNotice 删除公告
func (l *NoticeLogic) DeleteNotice(ctx context.Context, req types.DeleteNoticeReq) (*types.DeleteNoticeResp, error) {
	noticeRepo := repo.NewNoticeRepo(global.DB)
	if err := noticeRepo.DeleteNotice(req.ID); err != nil {
		zlog.CtxErrorf(ctx, "DeleteNotice failed: %v", err)
		return nil, response.ErrResp(err, codeNoticeDeleteFailed)
	}
	zlog.CtxInfof(ctx, "Notice deleted successfully (noticeID: %d)", req.ID)
	return &types.DeleteNoticeResp{}, nil
}

// GetNoticeByID 获取公告详情
func (l *NoticeLogic) GetNoticeByID(ctx context.Context, req types.GetNoticeReq) (*types.GetNoticeResp, error) {
	noticeRepo := repo.NewNoticeRepo(global.DB)
	notice, err := noticeRepo.GetNoticeByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "GetNoticeByID failed: notice not found (noticeID: %d)", req.ID)
			return nil, response.ErrResp(err, codeNoticeNotFound)
		}
		zlog.CtxErrorf(ctx, "GetNoticeByID failed: %v", err)
		return nil, response.ErrResp(err, response.INTERNAL_ERROR)
	}
	zlog.CtxInfof(ctx, "Notice retrieved successfully (noticeID: %d)", req.ID)
	return &types.GetNoticeResp{
		UpdateAt: notice.UpdatedAt.Unix(),
		Content:  notice.Content,
	}, nil
}
