package types

// @Title        wechat.go
// @Description
// @Create       XdpCs 2025-03-06 上午1:25
// @Update       XdpCs 2025-03-06 上午1:25
type WechatLoginReq struct {
	Code string `json:"code"`
}

type WechatLoginResp struct {
	Openid  string `json:"openid"`
	Unionid string `json:"unionid"`
}
