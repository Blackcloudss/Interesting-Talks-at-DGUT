package initalize

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/configs"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/model"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/pkg/database"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/pkg/mysqlx"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/pkg/redisx"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
)

// @Title        init_database.go
// @Description
// @Create       XdpCs 2025-02-16 上午12:56
// @Update       XdpCs 2025-02-16 上午12:56
func InitDataBase(config configs.Config) {
	switch config.DB.Driver {
	case "mysql":
		database.InitDataBases(mysqlx.NewMySql(), config)
	default:

		zlog.Fatalf("不支持的数据库驱动：%s", config.DB.Driver)
	}
	if config.App.Env != "pro" {
		err := global.DB.AutoMigrate()

		//迁移数据库所有的表
		migrateTables()

		if err != nil {
			zlog.Fatalf("数据库迁移失败！")
		}
	}
	zlog.Infof("数据库初始化成功！")
}
func InitRedis(config configs.Config) {
	if config.Redis.Enable {
		var err error
		global.Rdb, err = redisx.GetRedisClient(config)
		if err != nil {
			zlog.Errorf("无法初始化Redis : %v", err)
		}
	} else {
		zlog.Warnf("不使用Redis")
	}

}

// 自动迁移表 --通过GORM的AutoMigrate方法自动创建或更新数据库表结构 确保表结构与模型定义一致
func migrateTables() {
	global.DB.AutoMigrate(
		//自动迁移 用户展示表，确保表结构存在
		&model.UserDisplay{},
		//自动迁移 用户普通信息表，确保表结构存在
		&model.UserCommon{},
		//自动迁移 用户私有信息表，确保表结构存在
		&model.UserPrivate{},
		//自动迁移 用户授权表，确保表结构存在
		&model.UserAuth{},
		//自动迁移 消息表，确保表结构存在
		&model.Message{},
		//自动迁移 图片表，确保表结构存在
		&model.Image{},
		//自动迁移 博客表，确保表结构存在
		&model.Blog{},
		//自动迁移 点赞表，确保表结构存在
		&model.Like{},
		//自动迁移 关注表，确保表结构存在
		&model.Follow{},
		//自动迁移 收藏表，确保表结构存在
		&model.Collection{},
		//自动迁移 评论表，确保表结构存在
		&model.FirstComment{},
		//自动迁移 二级评论表，确保表结构存在
		&model.SecondComment{},
		//自动迁移 评论点赞表，确保表结构存在
		&model.CommentLike{},
		//自动迁移 通知表，确保表结构存在
		&model.Notice{},
		//自动迁移 搜索历史表，确保表结构存在
		&model.SearchHistory{},
	)
}
