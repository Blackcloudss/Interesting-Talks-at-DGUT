package repo

import (
	"fmt"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"gorm.io/gorm"
)

type NoticeRepo struct {
	DB *gorm.DB
}

func NewNoticeRepo(db *gorm.DB) *NoticeRepo {
	return &NoticeRepo{
		DB: db,
	}
}

// CreateNotice 创建公告
func (r *NoticeRepo) CreateNotice(req types.CreateNoticeReq) (notice model.Notice, err error) {
	err = r.DB.Create(&model.Notice{
		Content: req.Content,
	}).Error
	if err != nil {
		zlog.Errorf(fmt.Sprintf("创建公告失败: %v", err))
		return
	}
	return notice, err
}

// UpdateNotice 更新公告
func (r *NoticeRepo) UpdateNotice(req types.UpdateNoticeReq) error {
	err := r.DB.Model(&model.Notice{}).
		Where(fmt.Sprintf("%v = ?", ID), req.ID).Updates(&model.Notice{
		Content: req.Content,
	}).Error
	if err != nil {
		zlog.Errorf(fmt.Sprintf("更新公告失败: %v", err))
		return err
	}
	return nil
}

// DeleteNotice 删除公告
func (r *NoticeRepo) DeleteNotice(id int64) error {
	err := r.DB.Model(&model.Notice{}).
		Where(fmt.Sprintf("%v = ?", ID), id).
		Delete(&model.Notice{}).
		Error
	if err != nil {
		zlog.Errorf(fmt.Sprintf("删除公告失败: %v", err))
		return err
	}
	return nil
}

// GetNoticeByID 根据公告ID获取公告
func (r *NoticeRepo) GetNoticeByID(id int64) (*model.Notice, error) {
	var notice model.Notice
	err := r.DB.Model(&model.Notice{}).
		Where(fmt.Sprintf("%v = ?", ID), id).
		First(&notice).
		Error
	if err != nil {
		zlog.Errorf(fmt.Sprintf("获取公告失败: %v", err))
		return nil, err
	}
	return &notice, nil
}
