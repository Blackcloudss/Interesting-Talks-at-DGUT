package logic

import (
	"context"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
)

// @Title        wechat.go
// @Description
// @Create       XdpCs 2025-03-06 上午1:32
// @Update       XdpCs 2025-03-06 上午1:32

type WechatLogic struct {
}

func NewWechatLoginLogic() *WechatLogic {
	return &WechatLogic{}
}

func (l *WechatLogic) WechatLogin(ctx context.Context, req types.WechatLoginReq) {

}
