package logic

import (
	"context"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/jwt"
	"time"
)

// @Title        token.go
// @Description
// @Create       XdpCs 2025-03-10 下午3:51
// @Update       XdpCs 2025-03-10 下午3:51
type TokenLogic struct {
}

func NewTokenLogic() *TokenLogic {
	return &TokenLogic{}
}

// GenAtoken
//
//	@Description: 生成atoken
//	@receiver l
//	@param ctx
//	@param data
//	@return resp
//	@return err
func (l *TokenLogic) GenAtoken(ctx context.Context, data jwt.TokenData) (resp types.TokenResp, err error) {
	defer utils.RecordTime(time.Now())()
	resp.Atoken, err = jwt.GenToken(data)
	if err != nil {
		zlog.CtxErrorf(ctx, "AtokenLogic err:%v", err)
		return
	}
	return
}

// GenRtoken
//
//	@Description: 生成rtoken
//	@receiver l
//	@param ctx
//	@param data
//	@return resp
//	@return err
func (l *TokenLogic) GenRtoken(ctx context.Context, data jwt.TokenData) (resp types.TokenResp, err error) {
	defer utils.RecordTime(time.Now())()
	resp.Rtoken, err = jwt.GenToken(data)
	if err != nil {
		zlog.CtxErrorf(ctx, "TokenLogic err:%v", err)
	}
	data.Time = global.ATOKEN_EFFECTIVE_TIME
	data.Class = global.AUTH_ENUMS_ATOKEN
	resp.Atoken, err = jwt.GenToken(data)
	if err != nil {
		zlog.CtxErrorf(ctx, "RtokenLogic err:%v", err)
		return
	}
	return
}

// RefreshToken
//
//	@Description: 用rtoken刷新atoken和rtoken
//	@receiver l
//	@param ctx
//	@param req
//	@return resp
//	@return err
func (l *TokenLogic) RefreshToken(ctx context.Context, req types.TokenReq) (resp types.TokenResp, err error) {
	//解析token是否有效，并取出上一次的值
	data, err := jwt.IdentifyToken(ctx, req.Token)
	if err != nil {
		//对应token无效，直接让他返回
		return resp, response.ErrResp(err, response.TOKEN_IS_EXPIRED)
	}
	//判断其是否为rtoken
	if data.Class != global.AUTH_ENUMS_RTOKEN {
		return resp, response.ErrResp(err, response.TOKEN_TYPE_ERROR)
	}
	//生成新的token
	resp.Atoken, err = jwt.GenToken(jwt.FullToken(global.AUTH_ENUMS_ATOKEN, data.Userid))
	if err != nil {
		return resp, response.ErrResp(err, response.TOKEN_NOT_VALID)
	}
	resp.Rtoken, err = jwt.GenToken(jwt.TokenData{
		Class:  global.AUTH_ENUMS_RTOKEN,
		Time:   data.Time,
		Userid: data.Userid,
	})
	if err != nil {
		return resp, response.ErrResp(err, response.TOKEN_NOT_VALID)
	}
	return resp, nil
}
