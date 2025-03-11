package types

// @Title        login.go
// @Description
// @Create       XdpCs 2025-03-06 上午1:25
// @Update       XdpCs 2025-03-06 上午1:25

// 微信登录请求结构体
type WechatLoginReq struct {
	JsCode string `json:"js_code"` //微信临时登陆凭证
}

// 调用code2session接口后返回数据
type Code2SessionResp struct {
	Openid     string `json:"openid"`      // 用户在不同类型产品中的身份id，不同产品互不相同
	Unionid    string `json:"unionid"`     // 用户在微信开放平台的唯一标识符 不同产品都相同
	SessionKey string `json:"session_key"` // 会话密钥
	Errcode    int    `json:"errcode"`     // 错误码
	Errmsg     string `json:"errmsg"`      // 错误信息
}

// 微信登录响应结构体
type WechatLoginResp struct {
	Atoken string `json:"token"`  // 账号登录认证
	Rtoken string `json:"rtoken"` // 刷新处理
}
