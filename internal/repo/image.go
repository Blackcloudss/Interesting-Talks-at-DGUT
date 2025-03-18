package repo

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"gorm.io/gorm"
)

type ImageRepo struct {
	DB *gorm.DB
}

func NewImageRepo(db *gorm.DB) *ImageRepo {
	return &ImageRepo{DB: db}
}

// CreateImage 创建图片
func (r *ImageRepo) CreateImage(image *model.Image) error {
	return r.DB.Create(image).Error
}

//删除图片

//获取图片
