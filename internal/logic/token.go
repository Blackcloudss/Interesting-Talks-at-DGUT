package logic

import (
	"context"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/repo"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/jwt"
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
		return resp, response.ErrResp(err, response.TOKEN_NOT_VALID)
	}
	//判断其是否为rtoken
	if data.Class != global.AUTH_ENUMS_RTOKEN {
		return resp, response.ErrResp(err, response.TOKEN_TYPE_ERROR)
	}

	//生成新的token
	resp.Atoken, err = jwt.GenToken(jwt.FullToken(global.AUTH_ENUMS_ATOKEN, data.Userid))
	if err != nil {
		return resp, response.ErrResp(err, response.GET_ATOKEN_ERROR)
	}
	resp.Rtoken, err = jwt.GenToken(jwt.TokenData{
		Class:  global.AUTH_ENUMS_RTOKEN,
		Time:   data.Time,
		Userid: data.Userid,
	})
	if err != nil {
		return resp, response.ErrResp(err, response.GET_RTOKEN_ERROR)
	}
	//查找用户当前身份
	resp.Role, err = repo.NewUserRepo(global.DB).GetUserRole(data.Userid)
	if err != nil {
		return resp, response.ErrResp(err, GET_USER_ROLE)
	}
	return resp, nil
}
