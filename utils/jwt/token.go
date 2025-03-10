package jwt

import (
	"context"
	"errors"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// @Title        token.go
// @Description
// @Create       XdpCs 2025-02-24 下午11:06
// @Update       XdpCs 2025-02-24 下午11:06
type TokenData struct {
	Userid int64         //用户ID
	Class  string        //类别
	Time   time.Duration //过期时间
}

type MyClaims struct {
	Userid int64  `json:"userid"` //用户ID
	Type   string `json:"type"`   // 类型
	jwt.RegisteredClaims
}

var mySecret = []byte("InterestingTalk")

func GenToken(data TokenData) (string, error) {
	// 创建一个我们自己的声明
	claims := MyClaims{
		data.Userid,
		data.Class,
		jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(data.Time)), // 过期时间
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// 用于验证令牌是否有效
func IdentifyToken(ctx context.Context, Token string) (TokenData, error) {
	//解析token
	claim, err := ParseToken(Token)
	var data TokenData
	if err != nil {
		zlog.CtxErrorf(ctx, "IdentifyToken err: %v", err)
		return TokenData{}, err
	}
	data.Userid = claim.Userid
	data.Class = claim.Type
	if claim.Type == global.AUTH_ENUMS_RTOKEN {
		data.Time = global.RTOKEN_EFFECTIVE_TIME - time.Duration(time.Now().Unix()-claim.RegisteredClaims.NotBefore.Unix())
	} else {
		data.Time = global.ATOKEN_EFFECTIVE_TIME
	}
	return data, nil
}

func FullToken(class string, user_id int64) (data TokenData) {
	data.Userid = user_id
	if class == global.AUTH_ENUMS_ATOKEN {
		data.Time = global.ATOKEN_EFFECTIVE_TIME
		data.Class = global.AUTH_ENUMS_ATOKEN
	} else {
		data.Time = global.RTOKEN_EFFECTIVE_TIME
		data.Class = global.AUTH_ENUMS_RTOKEN
	}
	return
}

// GetUserId
//
//	@Description:
//	@param c
//	@return int64
func GetUserId(c *gin.Context) int64 {
	if data, exists := c.Get(global.TOKEN_USER_ID); exists {
		user_id, ok := data.(int64)
		if ok {
			return user_id
		}
	}
	return 0
}
