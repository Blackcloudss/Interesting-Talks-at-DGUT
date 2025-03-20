package logic

import (
	"bytes"
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
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/snowflake"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	REDIS_SET_FAULT = response.MsgCode{50001, "redis存取SessionKey失败"}
	GET_USER_ROLE   = response.MsgCode{50006, "获取用户角色失败"}
)

const (
	// 获取微信二维码
	GET_QRCODE_URL = "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s"
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
	client := &http.Client{Timeout: 5 * time.Second}
	result, err := client.Get(url)

	// 先检查错误再判断状态码
	if err != nil {
		zlog.CtxErrorf(ctx, "调用微信code2Session接口失败：%v", err)
		return resp, response.ErrResp(err, response.COMMON_FAIL)
	}

	// 后校验状态码
	if result.StatusCode != http.StatusOK {
		zlog.CtxErrorf(ctx, "微信接口异常，状态码：%d", result.StatusCode)
		return resp, response.ErrResp(err, response.COMMON_FAIL)
	}

	// 对 关闭Body 做封装处理
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			zlog.CtxErrorf(ctx, "关闭Body失败: %v", err)
		}
	}(result.Body)

	var C2S types.Code2SessionResp

	if err = json.NewDecoder(result.Body).Decode(&C2S); err != nil {
		zlog.CtxErrorf(ctx, "响应解析失败: %v", err)
		return resp, response.ErrResp(err, response.COMMON_FAIL)
	}

	//判断 该用户是否在数据库中,没有则存放数据库中
	UserId, err := repo.NewUserRepo(global.DB).JudgeUser(C2S.Openid)
	if err != nil {
		zlog.CtxErrorf(ctx, "GenLoginData err: %v", err)
		return resp, response.ErrResp(err, response.COMMON_FAIL)
	}
	//获取用户角色
	resp.Role, err = repo.NewUserRepo(global.DB).GetUserRole(UserId)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取用户角色失败: %v", err)
		return resp, response.ErrResp(err, GET_USER_ROLE)
	}

	//把用户新的Sessionkey放进Redis
	if err = global.Rdb.Set(ctx, fmt.Sprintf(global.REDIS_SESSIONKEY, C2S.Openid), C2S.SessionKey, global.SESSIONKEY_EFFECTIVE_TIME).Err(); err != nil {
		zlog.CtxErrorf(ctx, "redis set session_key err: %v", err)
		return resp, response.ErrResp(err, REDIS_SET_FAULT)
	}

	//制作 Atoken 和 Rtoken 自定义登陆态
	resp.Atoken, err = jwt.GenToken(jwt.FullToken(global.AUTH_ENUMS_ATOKEN, UserId))
	resp.Rtoken, err = jwt.GenToken(jwt.FullToken(global.AUTH_ENUMS_RTOKEN, UserId))

	return resp, nil
}

func (l *WechatLogic) GetQRCode(ctx context.Context) (resp *types.QRCodeResp, err error) {
	defer utils.RecordTime(time.Now())()

	// 获取微信access_token
	WxAtoken, err := jwt.JudgeWxAtoken(ctx)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取微信access_token失败：%v", err)
		return nil, response.ErrResp(err, response.GET_WXATOKEN_FAULT)
	}

	//用雪花算法生成随机且唯一的场景值
	scene := snowflake.GetString12Id(global.Node)
	req := types.QRCode{
		Page:      "pages/index/index", // 默认跳转到主页面
		Scene:     scene,
		Width:     280,
		IsHyaline: true,
		LineColor: types.LineColor{
			R: 0,
			G: 144,
			B: 0,
		},
	}

	Url := fmt.Sprintf(GET_QRCODE_URL, WxAtoken)

	jsonData, err := json.Marshal(req)
	if err != nil {
		zlog.CtxErrorf(ctx, "json.Marshal err: %v", err)
		return nil, err
	}

	result, err := http.Post(Url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		zlog.CtxErrorf(ctx, "http.Post err: %v", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			zlog.CtxErrorf(ctx, "关闭Body失败: %v", err)
		}
	}(result.Body)

	// 记录状态码
	resp = &types.QRCodeResp{
		Errcode: result.StatusCode,
	}

	// 统一处理响应
	if result.StatusCode != http.StatusOK {
		// 处理错误响应（JSON格式）
		if err = json.NewDecoder(result.Body).Decode(resp); err != nil {
			zlog.CtxErrorf(ctx, "错误响应解析失败: %v", err)
			return nil, response.ErrResp(err, response.COMMON_FAIL)
		}
		return resp, nil
	}

	// 处理成功响应（图片二进制）
	resp.Buffer, err = ioutil.ReadAll(result.Body)
	if err != nil {
		zlog.CtxErrorf(ctx, "图片数据读取失败: %v", err)
		return nil, response.ErrResp(err, response.COMMON_FAIL)
	}

	return
}
