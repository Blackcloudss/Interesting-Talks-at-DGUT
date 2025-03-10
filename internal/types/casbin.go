package types

// @Title        casbin.go
// @Description
// @Create       XdpCs 2025-03-11 上午12:37
// @Update       XdpCs 2025-03-11 上午12:37

// 获得权限组（入参）
type RuleReq struct {
	UserId int64 `form:"user_id"` // 仅为了测试使用，之后删除
}

// 获得权限组（出参）
type RuleResp struct {
	Url  []string `json:"urls"`  // 包含权限 URL 的数组
	Role string   `json:"level"` // 权限等级
}

// 权限验证
type RuleCheck struct {
	UserId int64 `form:"user_id"` // 测试时使用
}
