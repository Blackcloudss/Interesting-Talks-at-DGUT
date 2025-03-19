package logic

import (
	"context"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/repo"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils"
	"time"
)

var (
	GET_FRIENDS_ID   = response.MsgCode{51000, "获取好友们ID失败"}
	GET_FRIENDS_INFO = response.MsgCode{51001, "获取好友们信息失败"}
)

// @Title        chat.go
// @Description
// @Create       XdpCs 2025-03-19 下午3:07
// @Update       XdpCs 2025-03-19 下午3:07
type Chatlogic struct {
}

func NewChatlogic() *Chatlogic {
	return &Chatlogic{}
}

// GetFriendList
//
//	@Description: 获取好友列表
//	@receiver l
//	@param ctx
//	@param UserId
//	@return resp
//	@return err
func (l *Chatlogic) GetFriendList(ctx context.Context, UserId int64) (resp types.GetFriendListResp, err error) {
	defer utils.RecordTime(time.Now())()
	// 获取好友们id
	FriendsId, err := repo.NewChatRepo(global.DB).GetFriendsId(UserId)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取好友们id失败: %v", err)
		return resp, response.ErrResp(err, GET_FRIENDS_ID)
	}
	// 获取好友信息
	resp.Friends, err = repo.NewChatRepo(global.DB).GetFriendList(FriendsId)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取好友们信息失败: %v", err)
		return resp, response.ErrResp(err, GET_FRIENDS_INFO)
	}
	return resp, nil
}
