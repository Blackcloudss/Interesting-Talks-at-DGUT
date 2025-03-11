package types

// @Title        token.go
// @Description
// @Create       XdpCs 2025-03-06 上午11:18
// @Update       XdpCs 2025-03-06 上午11:18

// 后端接受前端传过来的token
type TokenReq struct {
	Token string `json:"token"`
}

// 后端将atoken和rtoken传回给前端
type TokenResp struct {
	Atoken string `json:"atoken"` //认证身份
	Rtoken string `json:"rtoken"` //刷新Atoken
}

// 后端请求微信的Atoken
type WxTokenReq struct {
	Grant_Type   string `json:"grant_type"`    // 权力类型 固定值  client_credential
	AppId        string `json:"appid"`         // 小程序的appid
	AppSecret    string `json:"secret"`        // 小程序的appsecret
	ForceRefresh bool   `json:"force_refresh"` //调用模式 false：普通调用模式 ，true:强制刷新模式
}

// 后端接受微信的token
type WxTokenResp struct {
	AccessToken string `json:"access_token"` // 身份凭证
	Expires_In  int    `json:"expires_in"`   // 有效期
}
