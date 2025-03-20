package repo

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/image"
	"gorm.io/gorm"
	"path/filepath"
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

func (r *ImageRepo) GetImagePathsByBlogID(blogID int64) ([]string, error) {
	var imagePaths []string
	if err := r.DB.Model(&model.Image{}).Where("blog_id = ?", blogID).Pluck("path", &imagePaths).Error; err != nil {
		return nil, err
	}
	return imagePaths, nil
}

// DeleteImagesByBlogID 根据blogid删除图片
func (r *ImageRepo) DeleteImagesByBlogID(blogID int64) error {
	// 获取所有需要删除的图片路径
	imagePaths, err := r.GetImagePathsByBlogID(blogID)
	if err != nil {
		zlog.Errorf("获取图片路径失败：%v", err)
		return err
	}

	// 删除本地文件
	for _, imagePath := range imagePaths {
		// 拼接完整的文件路径
		fullPath := filepath.Join("images", imagePath)
		if err := image.DeleteLocalFile(fullPath); err != nil {
			zlog.Errorf("删除本地文件失败：%v", err)
			return err
		}
	}

	// 删除数据库中的记录
	if err := r.DB.Where("blog_id = ?", blogID).Delete(&model.Image{}).Error; err != nil {
		zlog.Errorf("删除数据库记录失败：%v", err)
		return err
	}

	return nil
}
