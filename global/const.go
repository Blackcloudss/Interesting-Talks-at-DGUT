package global

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/filter_swords"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/snowflake"
	"time"
)

// @Title        url.go
// @Description
// @Create       XdpCs 2025-02-01 下午8:14
// @Update       XdpCs 2025-02-01 下午8:14

// 所有常量文件读取位置
const (
	DEFAULT_CONFIG_FILE_PATH = "/config.yaml"

	REDIS_SESSIONKEY          = "ITAD:openid.login:%s:string"
	SESSIONKEY_EFFECTIVE_TIME = time.Hour * 24 * 7
	REDIS_WXATOKEN_KEY        = "ITAD:WxAtoken:%s:string"
	WXATOKEN_EFFECTIVE_TIME   = time.Minute * 108
	REDIS_MESSAGES_KEY        = "Chat:%s, Type:%v"
	MESSAGES_EFFECTIVE_TIME   = time.Hour * 24 * 30

	ATOKEN_EFFECTIVE_TIME = time.Hour * 12
	RTOKEN_EFFECTIVE_TIME = time.Hour * 24 * 30

	AUTH_ENUMS_ATOKEN = "atoken"
	AUTH_ENUMS_RTOKEN = "rtoken"
	DEFAULT_NODE_ID   = 1
	TOKEN_USER_ID     = "UserId"

	MSGID        = "msg_id"
	STATUS       = "status"       // 状态
	MESSAGE      = "message"      // 消息
	ACK          = "ack"          // 客户端发送的确认请求类型（消息头标识）
	ACKNOWLEDGED = "acknowledged" // 服务端标记的最终确认状态 持久化状态
	OFFLINE      = "offline"      // 离线
	DELIVERED    = "delivered"    // 在线
)

var (
	Node, _ = snowflake.NewNode(DEFAULT_NODE_ID)

	TOURIST_URLS = []string{}

	STUDENT_URLS = []string{}

	MANAGER_URLS = []string{
		"/api/profile/role",  // 更改用户角色
		"/api/notice/create", // 创建公告
		"/api/notice/update", // 更新公告
		"/api/notice/delete", // 删除公告
		"/api/notice/show",   // 获取公告
	}

	Filter = filter_swords.New()
)
