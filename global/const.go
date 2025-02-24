package global

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/snowflake"
	"time"
)

// @Title        const.go
// @Description
// @Create       XdpCs 2025-02-01 下午8:14
// @Update       XdpCs 2025-02-01 下午8:14

// 所有常量文件读取位置
const (
	DEFAULT_CONFIG_FILE_PATH = "/config.yaml"
	REDIS_PHONE_CODE         = "Achobeta:phone.login.code:%s:string"
	REDIS_PHONE              = "Achobeta:phone.login.effective:%s:string"
	ATOKEN_EFFECTIVE_TIME    = time.Hour * 12
	RTOKEN_EFFECTIVE_TIME    = time.Hour * 24 * 30
	AUTH_ENUMS_ATOKEN        = "atoken"
	AUTH_ENUMS_RTOKEN        = "rtoken"
	DEFAULT_NODE_ID          = 1
	TOKEN_USER_ID            = "UserId"
)

var Node, _ = snowflake.NewNode(DEFAULT_NODE_ID)

var NORMAL_ADMIN_URLS = []string{}

var SUPER_ADMIN_URLS = []string{}
