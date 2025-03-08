package types

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// @Title        token.go
// @Description
// @Create       XdpCs 2025-03-06 上午11:18
// @Update       XdpCs 2025-03-06 上午11:18

type TokenData struct {
	Class  string
	Issuer string
	Time   time.Duration
	Userid int64
}

type MyClaims struct {
	Type   string `json:"type"`
	Userid int64  `json:"userid"`
	jwt.RegisteredClaims
}

type TokenReq struct {
	Token string `json:"token"`
}

type TokenResp struct {
	Atoken string `json:"atoken"`
	Rtoken string `json:"rtoken"`
}
