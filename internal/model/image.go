package model

// @Title        image.go
// @Description
// @Create       XdpCs 2025-03-17 上午2:21
// @Update       XdpCs 2025-03-17 上午2:21

type Image struct {
	CommonModel
	ImagePath string `gorm:"column:image_path;type:varchar(255);comment:'图片路径'"`
	Size      int64  `gorm:"column:size;type:bigint;comment:'文件大小'"`
	BlogID    int64  `gorm:"column:blog_id;type:bigint;comment:'博客ID'"`
}

func (Image) TableName() string { return "image" }
