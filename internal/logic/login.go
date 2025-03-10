package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/configs"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/repo"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/jwt"
	"net/http"
	"time"
)

// @Title        login.go
// @Description
// @Create       XdpCs 2025-03-06 上午1:32
// @Update       XdpCs 2025-03-06 上午1:32

type WechatLogic struct {
}

func NewWechatLoginLogic() *WechatLogic {
	return &WechatLogic{}
}

// 微信登陆
func (l *WechatLogic) WechatLogin(ctx context.Context, req types.WechatLoginReq) (resp *types.WechatLoginResp, err error) {
	defer utils.RecordTime(time.Now())()
	resp, err = WxLogin(ctx, req.JsCode)
	if err != nil {
		zlog.CtxErrorf(ctx, "调用微信code2Session接口失败：%v", err)
		return nil, response.ErrResp(err, response.COMMON_FAIL)
	}
	return resp, nil
}

// 向微信服务器请求code2Session
func WxLogin(ctx context.Context, code string) (resp *types.WechatLoginResp, err error) {
	url := fmt.Sprintf(configs.Conf.Wechat.BaseUrl, configs.Conf.Wechat.AppID, configs.Conf.Wechat.AppSecret, code)
	result, err := http.DefaultClient.Get(url)
	if err != nil {
		zlog.CtxErrorf(ctx, "调用微信code2Session接口失败：%v", err)
		return resp, response.ErrResp(err, response.COMMON_FAIL)
	}
	defer result.Body.Close()

	var C2S types.Code2SessionResp

	if err = json.NewDecoder(result.Body).Decode(&C2S); err != nil {
		zlog.CtxErrorf(ctx, "响应解析失败: %v", err)
		return nil, response.ErrResp(err, response.COMMON_FAIL)
	}

	//判断 该用户是否在数据库中,没有则存放数据库中
	UserId, err := repo.NewUserRepo(global.DB).JudgeUser(C2S.Openid)
	if err != nil {
		zlog.CtxErrorf(ctx, "GenLoginData err: %v", err)
		return resp, response.ErrResp(err, response.COMMON_FAIL)
	}

	//把用户的Sessionkey放进Redis
	if err = global.Rdb.Set(ctx, fmt.Sprintf(global.REDIS_SESSION_KEY, C2S.Openid), C2S.SessionKey, global.SESSIONKEY_EFFECTIVE_TIME).Err(); err != nil {
		zlog.CtxErrorf(ctx, "redis set session_key err: %v", err)
		return resp, response.ErrResp(err, response.COMMON_FAIL)
	}

	//制作成 Atoken 和 Rtoken 登陆态
	resp.Atoken, err = jwt.GenToken(jwt.FullToken(global.AUTH_ENUMS_ATOKEN, UserId))
	resp.Rtoken, err = jwt.GenToken(jwt.FullToken(global.AUTH_ENUMS_RTOKEN, UserId))

	return resp, nil
}
