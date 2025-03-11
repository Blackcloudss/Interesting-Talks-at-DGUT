package logic

import (
	"context"
	"encoding/json"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/configs"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	BLANK_TOKEN = ""
	ATOKEN_URL  = "https://api.weixin.qq.com/cgi-bin/stable_token"
	PHONE_URL   = "https://api.weixin.qq.com/wxa/business/getuserphonenumber"
	GRANT_TYPE  = "client_credential"
)

var (
	GET_WXATOKEN_FAULT = response.MsgCode{50002, "获取微信的access_token失败"}
	GET_PHONE_FAULT    = response.MsgCode{50003, "获取微信手机号失败"}
)

// @Title        phone.go
// @Description
// @Create       XdpCs 2025-03-11 上午9:33
// @Update       XdpCs 2025-03-11 上午9:33
type Phonelogic struct {
}

func NewPhoneLogic() *Phonelogic {
	return &Phonelogic{}
}

// 获取用户手机号
func (l *Phonelogic) GetPhone(ctx context.Context, req types.WxPhoneReq) (resp *types.WxPhoneResp, err error) {
	defer utils.RecordTime(time.Now())()
	//获取微信的access_token
	req.WxAtoken, err = l.GetWxAtoken(ctx)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取微信的access_token失败：%v", err)
		return nil, response.ErrResp(err, GET_WXATOKEN_FAULT)
	}

	//调用微信的getPhoneNumber接口
	result, err := l.GetPhoneNumber(ctx, req.WxAtoken, req.Code)
	if err != nil {
		zlog.CtxErrorf(ctx, "调用微信getPhoneNumber接口失败：%v", err)
		return nil, response.ErrResp(err, GET_PHONE_FAULT)
	}
	resp = new(types.WxPhoneResp)
	resp.PhoneNumber = result.PhoneInfo.PhoneNumber
	resp.PurePhoneNumber = result.PhoneInfo.PurePhoneNumber
	return
}

// GetWxAtoken
//
//	@Description:
//	@receiver l
//	@param ctx
//	@return WxAtoken
//	@return err
//
// 获取微信的access_token, 每次调用不强制刷新
func (l *Phonelogic) GetWxAtoken(ctx context.Context) (WxAtoken string, err error) {
	client := &http.Client{Timeout: 5 * time.Second}
	params := url.Values{}
	params.Add("grant_type", GRANT_TYPE)
	params.Add("appid", configs.Conf.Wechat.AppID)
	params.Add("secret", configs.Conf.Wechat.AppSecret)

	result, err := client.PostForm(ATOKEN_URL, params)
	// 先检查错误再判断状态码
	if err != nil {
		zlog.CtxErrorf(ctx, "调用微信getStableAccessToken接口失败：%v", err)
		return BLANK_TOKEN, response.ErrResp(err, response.COMMON_FAIL)
	}

	// 后校验状态码
	if result.StatusCode != http.StatusOK {
		zlog.CtxErrorf(ctx, "微信接口异常，状态码：%d", result.StatusCode)
		return BLANK_TOKEN, response.ErrResp(err, response.COMMON_FAIL)
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
		zlog.CtxErrorf(ctx, "获取的Atoken为空值：%v", err)
		return BLANK_TOKEN, response.ErrResp(err, response.COMMON_FAIL)
	}
	WxAtoken = req.AccessToken
	return
}

// GetPhoneNumber
//
//	@Description:
//	@receiver l
//	@param ctx
//	@param atoken
//	@param code
//	@return resp
//	@return err
//
// 调用微信的getPhoneNumber接口 获取手机号
func (l *Phonelogic) GetPhoneNumber(ctx context.Context, atoken, code string) (resp *types.WxPhone, err error) {
	client := &http.Client{Timeout: 5 * time.Second}
	params := url.Values{}
	params.Add("access_token", atoken)
	params.Add("code", code)

	result, err := client.PostForm(PHONE_URL, params)
	// 先检查错误再判断状态码
	if err != nil {
		zlog.CtxErrorf(ctx, "调用微信getPhoneNumber接口失败：%v", err)
		return nil, response.ErrResp(err, response.COMMON_FAIL)
	}

	// 后校验状态码
	if result.StatusCode != http.StatusOK {
		zlog.CtxErrorf(ctx, "微信接口异常，状态码：%d", result.StatusCode)
		return nil, response.ErrResp(err, response.COMMON_FAIL)
	}
	// 对 关闭Body 做封装处理
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			zlog.CtxErrorf(ctx, "关闭Body失败: %v", err)
			return
		}
	}(result.Body)

	//初始化结构体
	resp = new(types.WxPhone)
	if err = json.NewDecoder(result.Body).Decode(resp); err != nil {
		zlog.CtxErrorf(ctx, "响应解析失败: %v", err)
		return nil, response.ErrResp(err, response.COMMON_FAIL)
	}
	if resp.Errcode != 0 {
		zlog.CtxErrorf(ctx, "微信接口异常，状态码：%d", resp.Errcode)
		zlog.CtxErrorf(ctx, "微信接口异常，错误信息：%s", resp.Errmsg)
		return nil, response.ErrResp(err, response.COMMON_FAIL)
	}

	return
}
