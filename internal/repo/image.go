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

/*
// DeleteImage 根据图片ID删除图片
func (r *ImageRepo) DeleteImage(imageID int) error {
	return r.DB.Delete(&model.Image{}, imageID).Error
}*/

// GetImagesByBlogID 根据blogid获取图片
func (r *ImageRepo) GetImagesByBlogID(blogID int64) ([]*model.Image, error) {
	var images []*model.Image
	result := r.DB.Where("blog_id = ?", blogID).Find(&images)
	return images, result.Error
}
