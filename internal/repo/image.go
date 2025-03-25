package repo

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/image"
	"gorm.io/gorm"
	"os"
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

// GetImagesByBlogID 根据blogid获取图片
func (r *ImageRepo) GetImagesByBlogID(blogID int64) ([]*model.Image, error) {
	var images []*model.Image
	result := r.DB.Where("blog_id = ?", blogID).Find(&images)
	return images, result.Error
}

// GetImagePathsByBlogID 根据blogid获取图片路径
func (r *ImageRepo) GetImagePathsByBlogID(blogID int64) ([]string, error) {
	var imagePaths []string
	if err := r.DB.Model(&model.Image{}).Where("blog_id = ?", blogID).Pluck("image_path", &imagePaths).Error; err != nil {
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
		// 检查路径是否为文件
		fileInfo, err := os.Stat(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				// 文件不存在，跳过
				continue
			}
			zlog.Errorf("检查文件是否存在失败：%v", err)
			return err
		}
		if !fileInfo.IsDir() {
			// 确保是文件才删除
			if err := image.DeleteLocalFile(fullPath); err != nil {
				zlog.Errorf("删除本地文件失败：%v", err)
				return err
			}
		}
	}

	// 删除数据库中的记录
	if err := r.DB.Where("blog_id = ?", blogID).Delete(&model.Image{}).Error; err != nil {
		zlog.Errorf("删除数据库记录失败：%v", err)
		return err
	}

	return nil
}
