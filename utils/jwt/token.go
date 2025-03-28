package jwt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/configs"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
)

const (
	BLANK_TOKEN = ""
	ATOKEN_URL  = "https://api.weixin.qq.com/cgi-bin/stable_token"
	GRANT_TYPE  = "client_credential"
)

var (
	USE_WXAPI_ERROR   = response.MsgCode{50001, "调用微信getStableAccessToken接口失败"}
	WXAPI_ERROR       = response.MsgCode{50002, "微信接口异常"}
	WXATOKEN_IS_BLANK = response.MsgCode{50003, "获取的微信access_token是空值"}
	REDIS_SET_FAULT   = response.MsgCode{50001, "redis存取微信access_token失败"}
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
		// 计算剩余时间
		data.Time = global.RTOKEN_EFFECTIVE_TIME - time.Duration(time.Now().Unix()-claim.RegisteredClaims.NotBefore.Unix())
	} else {
		// 计算剩余时间
		data.Time = global.ATOKEN_EFFECTIVE_TIME
	}
	return data, nil
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

// 判断微信Atoken是否存在
func JudgeWxAtoken(ctx context.Context) (WxAtoken string, err error) {
	exists, err := global.Rdb.Exists(ctx, global.REDIS_WXATOKEN_KEY).Result()
	if err != nil && err != redis.Nil {
		zlog.CtxErrorf(ctx, "Redis访问异常: %v", err)
		return BLANK_TOKEN, response.ErrResp(errors.New("缓存服务异常"), response.GET_WXATOKEN_FAULT)
	}
	if exists == 0 {
		// 双检锁降低并发场景下的重复获取概率, 避免缓存击穿(避免大量客户端同时触发续期)
		var mu sync.Mutex
		mu.Lock()
		defer mu.Unlock()

		// 再次检查防止锁内已更新
		exists, err = global.Rdb.Exists(ctx, global.REDIS_WXATOKEN_KEY).Result()
		if exists == 0 {
			WxAtoken, err = GetWechatAccessToken(ctx)
			if err != nil {
				zlog.CtxErrorf(ctx, "获取微信access_token失败: %v", err)
				return BLANK_TOKEN, err
			}

			// 设置缓存并保留10%的冗余时间,防止缓存与真实Token同时失效
			if err = global.Rdb.Set(ctx, global.REDIS_WXATOKEN_KEY, WxAtoken, global.WXATOKEN_EFFECTIVE_TIME).Err(); err != nil {
				zlog.CtxErrorf(ctx, "缓存写入失败: %v", err)
				return BLANK_TOKEN, response.ErrResp(errors.New("缓存更新失败"), response.GET_WXATOKEN_FAULT)
			}
		}
	}
	return WxAtoken, nil
}

// GetWechatAccessToken
//
//	@Description:
//	@receiver l
//	@param ctx
//	@return WxAtoken
//	@return err
//
// 获取微信的access_token, 每次调用不强制刷新
func GetWechatAccessToken(ctx context.Context) (string, error) {
	// 1. 先检查配置是否正确
	if global.Config.Wechat.AppID == "" || global.Config.Wechat.AppSecret == "" {
		return "", errors.New("微信配置参数缺失")
	}

	// 检查配置文件中的参数
	fmt.Printf("AppID: %s\n", global.Config.Wechat.AppID)
	fmt.Printf("AppSecret: %s\n", global.Config.Wechat.AppSecret)

	// 2. 添加重试机制
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		token, err := tryGetAccessToken(ctx)
		if err == nil && token != "" {
			return token, nil
		}

		// 如果失败，等待一段时间后重试
		time.Sleep(time.Second * 2)
	}

	return "", errors.New("无法获取微信access token")
}

func tryGetAccessToken(ctx context.Context) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	params := url.Values{}
	params.Add("grant_type", GRANT_TYPE)
	params.Add("appid", configs.Conf.Wechat.AppID)
	params.Add("secret", configs.Conf.Wechat.AppSecret)

	result, err := client.PostForm(ATOKEN_URL, params)
	// 先检查错误再判断状态码
	if err != nil {
		zlog.CtxErrorf(ctx, "调用微信getStableAccessToken接口失败：%v", err)
		return BLANK_TOKEN, response.ErrResp(err, USE_WXAPI_ERROR)
	}

	// 后校验状态码
	if result.StatusCode != http.StatusOK {
		zlog.CtxErrorf(ctx, "微信接口异常，状态码：%d", result.StatusCode)
		return BLANK_TOKEN, response.ErrResp(err, WXAPI_ERROR)
	}
	// 对 关闭Body 做封装处理
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			zlog.CtxErrorf(ctx, "关闭Body失败: %v", err)
			return
		}
	}(result.Body)

	var req types.WxTokenResp

	if err = json.NewDecoder(result.Body).Decode(&req); err != nil {
		zlog.CtxErrorf(ctx, "响应解析失败: %v", err)
		return BLANK_TOKEN, response.ErrResp(err, response.COMMON_FAIL)
	}
	if req.AccessToken == "" {
		zlog.CtxErrorf(ctx, "获取的微信atoken为空值：%v", err)
		return BLANK_TOKEN, response.ErrResp(err, WXATOKEN_IS_BLANK)
	}
	WxAtoken := req.AccessToken
	zlog.CtxInfof(ctx, "获取微信access_token成功: %v", WxAtoken)
	// 将微信的atoken存储到redis中
	key := fmt.Sprintf(global.REDIS_WXATOKEN_KEY, configs.Conf.Wechat.AppID)
	if err = global.Rdb.Set(ctx, key, WxAtoken, global.WXATOKEN_EFFECTIVE_TIME).Err(); err != nil {
		zlog.CtxErrorf(ctx, "redis 存储微信Atoken失败: %v", err)
		return BLANK_TOKEN, response.ErrResp(err, REDIS_SET_FAULT)
	}
	return WxAtoken, nil
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
