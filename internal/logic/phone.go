package logic

import (
	"context"
	"encoding/json"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/repo"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/jwt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	PHONE_URL = "https://api.weixin.qq.com/wxa/business/getuserphonenumber"
)

var (
	GET_PHONE_FAULT  = response.MsgCode{50003, "获取微信手机号失败"}
	SAVE_PHONE_FAULT = response.MsgCode{50004, "保存手机号失败"}
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
func (l *Phonelogic) GetPhone(ctx context.Context, req types.WxPhoneReq, UserId int64) (resp *types.WxPhoneResp, err error) {
	defer utils.RecordTime(time.Now())()

	//检验redis中是否有Wxaccess_token,如果没有，需要重新获取
	WxAtoken, err := jwt.JudgeWxAtoken(ctx)

	//调用微信的getPhoneNumber接口
	result, err := l.GetPhoneNumber(ctx, WxAtoken, req.Code)
	if err != nil {
		zlog.CtxErrorf(ctx, "调用微信getPhoneNumber接口失败：%v", err)
		return nil, response.ErrResp(err, GET_PHONE_FAULT)
	}
	resp = new(types.WxPhoneResp)
	if result.PhoneInfo != nil {
		resp.PurePhoneNumber = result.PhoneInfo.PurePhoneNumber
	} else {
		zlog.CtxErrorf(ctx, "微信接口未返回 phone_info")
		return nil, response.ErrResp(err, GET_PHONE_FAULT)
	}
	//将手机号保存到数据库中
	err = repo.NewPhoneRepo(global.DB).SavePhone(UserId, resp.PurePhoneNumber)
	if err != nil {
		zlog.CtxErrorf(ctx, "保存手机号到数据库失败：%v", err)
		return nil, response.ErrResp(err, SAVE_PHONE_FAULT)
	}
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
		return nil, response.ErrResp(err, GET_PHONE_FAULT)
	}

	// 后校验状态码
	if result.StatusCode != http.StatusOK {
		zlog.CtxErrorf(ctx, "微信接口异常，状态码：%d", result.StatusCode)
		return nil, response.ErrResp(err, GET_PHONE_FAULT)
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
		return nil, response.ErrResp(err, GET_PHONE_FAULT)
	}
	if resp.Errcode != 0 {
		zlog.CtxErrorf(ctx, "微信接口异常，状态码：%d", resp.Errcode)
		zlog.CtxErrorf(ctx, "微信接口异常，错误信息：%s", resp.Errmsg)
		return nil, response.ErrResp(err, GET_PHONE_FAULT)
	}

	return
}
