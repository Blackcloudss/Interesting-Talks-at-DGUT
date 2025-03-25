package repo

import (
	"errors"
	"fmt"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"gorm.io/gorm"
)

const (
	LEVELONE   = "大一"
	LEVELTWO   = "大二"
	LEVELTHREE = "大三"
	LEVELFOUR  = "大四"
	LEVELFIVE  = "研一"
	LEVELSIX   = "研二"
	LEVELSEVEN = "新人"
	LEVELEIGHT = "实习生"
	LEVELNINE  = "元老"

	LEVELUPONE   = 10
	LEVELUPTWO   = 30
	LEVELUPTHREE = 50
	LEVELUPFOUR  = 100
	LEVELUPFIVE  = 150
	LEVELUPSIX   = 200
	LEVELUPSEVEN = 300
	LEVELUPEIGHT = 500
)

// @Title        point.go
// @Description
// @Create       XdpCs 2025-03-20 上午11:00
// @Update       XdpCs 2025-03-20 上午11:00
type PointRepo struct {
	DB *gorm.DB
}

func NewPointRepo(db *gorm.DB) *PointRepo {
	return &PointRepo{DB: db}
}

// 用户积分状态
type Status struct {
	TotalPoint int    `json:"total_point"`
	Tag        string `json:"tag"`
}

// 定义等级阈值列表（按阈值降序排列）
type levelThreshold struct {
	threshold int
	tag       string
}

func (r *PointRepo) IncreasePoints(userId int64, point int) (err error) {
	// 判断积分值是否合法
	if point <= 0 {
		return errors.New("invalid points value")
	}
	var status Status
	// 获取用户当前总积分
	err = r.DB.Model(&model.UserDisplay{}).
		Where(fmt.Sprintf("%s = ?", ID), userId).
		First(&status).
		Error
	if err != nil {
		zlog.Errorf("查询用户状态失败：%v", err)
		return err
	}

	TotalPoint := status.TotalPoint + point
	if TotalPoint > 10000 {
		zlog.Errorf("用户积分超过10000，无法继续增加积分")
		return nil
	}
	if status.Tag == LEVELNINE {
		zlog.Errorf("用户已达到最大等级，无法继续增加积分")
		return nil
	}
	// 判断用户等级
	var levelThresholds = []levelThreshold{
		{LEVELUPEIGHT, LEVELNINE},
		{LEVELUPSEVEN, LEVELEIGHT},
		{LEVELUPSIX, LEVELSEVEN},
		{LEVELUPFIVE, LEVELSIX},
		{LEVELUPFOUR, LEVELFIVE},
		{LEVELUPTHREE, LEVELFOUR},
		{LEVELUPTWO, LEVELTHREE},
		{LEVELUPONE, LEVELTWO},
		{0, LEVELONE}, // 兜底默认等级
	}

	// 动态匹配等级
	for _, lt := range levelThresholds {
		if TotalPoint >= lt.threshold {
			status.Tag = lt.tag
			break
		}
	}

	// 更新用户总积分和等级
	err = r.DB.Model(&model.UserDisplay{}).
		Where(fmt.Sprintf("%s = ?", ID), userId).
		Updates(&model.UserDisplay{
			TotalPoint: TotalPoint,
			Tag:        status.Tag,
		}).Error
	if err != nil {
		zlog.Errorf("更新用户总积分和等级失败：%v", err)
		return err
	}
	return nil
}
