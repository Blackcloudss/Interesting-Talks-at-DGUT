package logic

import (
	"context"
	"errors"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/repo"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils"
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
	codeUserFoundField = response.MsgCode{Code: 40013, Msg: "用户查询失败"}
	codeUserNotFound   = response.MsgCode{Code: 40014, Msg: "用户不存在"}
)

// TestLogic 逻辑层 用做逻辑处理相关操作
func (l *TestLogic) TestLogic(ctx context.Context, req types.TestO1Req) (resp *types.Test01Resp, err error) {
	defer utils.RecordTime(time.Now())()
	//..... some logic

	user, err := repo.NewTestRepo().GetUserById(req.UserID)
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
