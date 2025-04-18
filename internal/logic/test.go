package logic

import (
	"context"
	"errors"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/repo"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/jwt"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/snowflake"
	"gorm.io/gorm"
	"time"
)

// @Title        test.go
// @Description
// @Create       XdpCs 2025-02-24 下午11:30
// @Update       XdpCs 2025-02-24 下午11:30
type TestLogic struct {
}

func NewTestLogic() *TestLogic {
	return &TestLogic{}
}

// 这里定义我们内部logic的错误（非公有的常见类型的错误）
var (
	codeUserFoundField         = response.MsgCode{Code: 40013, Msg: "用户查询失败"}
	codeUserNotFound           = response.MsgCode{Code: 40014, Msg: "用户不存在"}
	codeBlogNotFound           = response.MsgCode{Code: 40020, Msg: "帖子不存在"}
	codeTransactionFailed      = response.MsgCode{Code: 40025, Msg: "事务处理失败"}
	codeInsufficientPermission = response.MsgCode{Code: 40034, Msg: "权限不足"}
)

// TestLogic 逻辑层 用做逻辑处理相关操作
func (l *TestLogic) TestLogic(ctx context.Context, req types.TestO1Req) (resp *types.Test01Resp, err error) {
	defer utils.RecordTime(time.Now())()
	//..... some logic
	user, err := repo.NewTestRepo(global.DB).GetUserById(req.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "user not found: %v", err)
			return nil, response.ErrResp(err, codeUserNotFound)
		} else {
			zlog.Errorf("get user error: %v", err)
			//注意 codeUserFoundField 只是事例，具体根据实际情况定义 response.ErrResp这里用作包装错误和响应，使得错误更加通用
			return nil, response.ErrResp(err, codeUserFoundField)
		}

	}
	resp.Name = user.Name
	resp.Age = user.Age

	return
}

func (l *TestLogic) CreateMember(ctx context.Context) (resp *types.CreateMemberResp, err error) {
	defer utils.RecordTime(time.Now())()
	resp = &types.CreateMemberResp{}
	Testid := snowflake.GetString12Id(global.Node)
	//判断 该用户是否在数据库中,没有则存放数据库中
	resp.UserId, err = repo.NewUserRepo(global.DB).JudgeUser(Testid)
	if err != nil {
		zlog.CtxErrorf(ctx, "GenLoginData err: %v", err)
		return resp, response.ErrResp(err, response.COMMON_FAIL)
	}

	//获取用户角色
	resp.Role, err = repo.NewUserRepo(global.DB).GetUserRole(resp.UserId)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取用户角色失败: %v", err)
		return resp, response.ErrResp(err, GET_USER_ROLE)
	}

	//制作 Atoken 和 Rtoken 自定义登陆态
	resp.Atoken, err = jwt.GenToken(jwt.FullToken(global.AUTH_ENUMS_ATOKEN, resp.UserId))
	if err != nil {
		return resp, response.ErrResp(err, response.GET_ATOKEN_ERROR)
	}
	resp.Rtoken, err = jwt.GenToken(jwt.FullToken(global.AUTH_ENUMS_RTOKEN, resp.UserId))
	if err != nil {
		return resp, response.ErrResp(err, response.GET_RTOKEN_ERROR)
	}
	return
}
