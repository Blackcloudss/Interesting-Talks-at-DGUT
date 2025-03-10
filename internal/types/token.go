package types

// @Title        token.go
// @Description
// @Create       XdpCs 2025-03-06 上午11:18
// @Update       XdpCs 2025-03-06 上午11:18

type TokenReq struct {
	Token string `json:"token"`
}

type TokenResp struct {
	Atoken string `json:"atoken"` //认证身份
	Rtoken string `json:"rtoken"` //刷新Atoken
}
