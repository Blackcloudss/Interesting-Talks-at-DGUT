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
	Role   string `json:"role"`   // 用户身份
}

// 生成小程序码的请求参数
type QRCodeReq struct {
	Scene     string    `json:"scene"`      // 必填，场景值（长度≤32字符）
	Page      string    `json:"page"`       // 可选，不填默认跳转主页
	Width     int       `json:"width"`      // 可选，二维码宽度（默认430px）
	LineColor LineColor `json:"line_color"` // 可选，使用 rgb 设置颜色
	IsHyaline bool      `json:"is_hyaline"` // 可选，是否需要透明底色，为 true 时，生成透明底色的小程序码
}

// 二维码颜色设置
type LineColor struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

// 生成小程序码的响应参数
type QRCodeResp struct {
	Errcode int    `json:"errcode"` // 错误码
	Errmsg  string `json:"errmsg"`  // 错误信息
	Buffer  []byte `json:"buffer"`  // 二维码图片
}
