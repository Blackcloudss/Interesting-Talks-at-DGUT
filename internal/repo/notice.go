package repo

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
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
	err = r.DB.Transaction(func(tx *gorm.DB) error {
		notice := model.Notice{
			Title:   req.Title,
			Content: req.Content,
			UserID:  req.UserID,
		}
		if err := tx.Create(&notice).Error; err != nil {
			return err
		}
		return nil
	})
	return notice, err
}

// UpdateNotice 更新公告
func (r *NoticeRepo) UpdateNotice(req types.UpdateNoticeReq) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Notice{}).Where("id = ?", req.ID).Updates(model.Notice{
			Title:   req.Title,
			Content: req.Content,
		}).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeleteNotice 删除公告
func (r *NoticeRepo) DeleteNotice(id int64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.Notice{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetNoticeByID 根据公告ID获取公告
func (r *NoticeRepo) GetNoticeByID(id int64) (*model.Notice, error) {
	var notice model.Notice
	if err := r.DB.First(&notice, id).Error; err != nil {
		return nil, err
	}
	return &notice, nil
}
