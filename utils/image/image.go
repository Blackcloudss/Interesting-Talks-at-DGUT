package image

import (
	"fmt"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/snowflake"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	IMAGE_PATH = "images"
	IMAGE_SIZE = 1024 * 1024 * 5
)

// @Title        image.go
// @Description
// @Create       XdpCs 2025-03-17 下午3:09
// @Update       XdpCs 2025-03-17 下午3:09
func UploadImage(file *multipart.FileHeader) (ImageUrl string, err error) {
	//图片大小限制
	if file.Size > IMAGE_SIZE {
		zlog.Errorf("图片大小超过限制")
		return
	}
	//图片格式限制
	filename := file.Filename
	ext, exist := ImageExtJudge(filename)
	if !exist {
		zlog.Errorf("图片格式不正确")
		return
	}
	// 生成唯一文件名
	newFilename := snowflake.GetString12Id(global.Node) + ext

	// 按日期分目录存储
	dateDir := time.Now().Format("2006-01") // 格式如 "2023-10"
	saveDir := filepath.Join(IMAGE_PATH, dateDir)
	savePath := filepath.Join(saveDir, newFilename)

	// 创建目录
	if err = os.MkdirAll(saveDir, os.ModePerm); err != nil {
		zlog.Errorf("创建目录失败:%v", err)
		return
	}

	// 保存文件到本地
	if err = SaveUploadFile(file, savePath); err != nil {
		zlog.Errorf("保存文件失败:%v", err)
		return "", err
	}

	// 返回相对路径（如 /images/2023-10/雪花id.jpg）
	ImageUrl = fmt.Sprintf("/%s/%s/%s", IMAGE_PATH, dateDir, newFilename)
	return
}

// 校验文件类型和大小
var whiteList = []string{".jpg", ".png", ".jpeg", ".gif", ".webp"}

func ImageExtJudge(FileName string) (ext string, exist bool) {
	ext = filepath.Ext(FileName)
	for _, v := range whiteList {
		if strings.ToLower(ext) == v {
			return ext, true
		}
	}
	return "", false
}

// 保存文件到本地
func SaveUploadFile(file *multipart.FileHeader, savePath string) error {
	// 打开文件
	src, err := file.Open()
	if err != nil {
		zlog.Errorf("打开文件失败:%v", err)
		return err
	}
	defer src.Close()

	// 创建本地文件
	out, err := os.Create(savePath)
	if err != nil {
		zlog.Errorf("创建本地文件失败:%v", err)
		return err
	}
	defer out.Close()

	// 将文件内容写入到本地文件
	if _, err = io.Copy(out, src); err != nil {
		zlog.Errorf("写入本地文件失败:%v", err)
		return err
	}
	return nil
}

// 删除本地文件
func DeleteLocalFile(filePath string) error {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		zlog.Warnf("文件不存在，无法删除：%s", filePath)
		return nil // 如果文件不存在，直接返回
	}
	// 删除文件
	if err := os.Remove(filePath); err != nil {
		zlog.Errorf("删除文件失败：%v", err)
		return err
	}
	return nil
}
