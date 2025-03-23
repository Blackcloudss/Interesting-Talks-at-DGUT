package api

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/jwt"
	"github.com/gin-gonic/gin"
)

// @Title        friend.go
// @Description
// @Create       XdpCs 2025-03-20 上午12:28
// @Update       XdpCs 2025-03-20 上午12:28

// GetFriendList
//
//	@Description: 获取好友列表
//	@param c
func GetFriendList(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	UserID := jwt.GetUserId(c)
	req, err := types.BindReq[types.GetFriendListReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetFriendList request error: %v", err)
		return
	}
	zlog.CtxInfof(ctx, "GetFriendList request: %v", req)
	resp, err := logic.NewFriendlogic().GetFriendList(ctx, UserID)
	response.Response(c, resp, err)
	return
}
