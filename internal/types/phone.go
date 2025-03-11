package types

// @Title        phone.go
// @Description
// @Create       XdpCs 2025-03-11 上午9:12
// @Update       XdpCs 2025-03-11 上午9:12

// 调用微信getPhoneNumber接口获取手机号
type WxPhoneReq struct {
	WxAtoken string `json:"atoken"` // 微信开发平台的Atoken
	Code     string `json:"code"`   //动态令牌   --与js_code不同,二者不能混用
}

// 将手机号传回给前端
type WxPhoneResp struct {
	PhoneNumber     string `json:"phoneNumber"`     //用户绑定的手机号（国外手机号会有区号)
	PurePhoneNumber string `json:"purePhoneNumber"` //没有区号的手机号
}

// 具体手机号信息
type WxPhone struct {
	Errcode   int       `json:"errcode"`    // 错误码
	Errmsg    string    `json:"errmsg"`     // 错误信息
	PhoneInfo PhoneInfo `json:"phone_info"` // 手机号信息
}

type PhoneInfo struct {
	PhoneNumber     string `json:"phoneNumber"`     //用户绑定的手机号（国外手机号会有区号)
	PurePhoneNumber string `json:"purePhoneNumber"` //没有区号的手机号
	CountryCode     string `json:"countryCode"`     //区号
	Watermark       struct {
		Timestamp int64  `json:"timestamp"` // 时间戳
		Appid     string `json:"appid"`     // 小程序appid
	} `json:"watermark"` // 水印信息
}
