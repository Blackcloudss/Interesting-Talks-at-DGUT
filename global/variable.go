package global

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/configs"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// @Title        variable.go
// @Description
// @Create       XdpCs 2025-02-01 下午8:15
// @Update       XdpCs 2025-02-01 下午8:15
var (
	Path   string
	DB     *gorm.DB
	Rdb    *redis.Client
	Config *configs.Config
)
