package api

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/jwt"
	"github.com/gin-gonic/gin"
)

// @Title        chat.go
// @Description
// @Create       XdpCs 2025-03-19 下午3:04
// @Update       XdpCs 2025-03-19 下午3:04
//
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
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "GetFriendList request: %v", req)
	resp, err := logic.NewChatlogic().GetFriendList(ctx, UserID)
	response.Response(c, resp, err)
	return
}

// GetFriendList
//
//	@Description: 加载历史消息
//	@param c
func GetMessageHistory(c *gin.Context) {
	//ctx := zlog.GetCtxFromGin(c)
	//req, err := types.BindReq[types.Tes](c)
	//if err != nil {
	//	zlog.CtxErrorf(ctx, "GetMessageHistory request error: %v", err)
	//	response.NewResponse(c).Error(response.PARAM_NOT_VALID)
	//	return
	//}
	//zlog.CtxInfof(ctx, "GetMessageHistory request: %v", req)
	//resp, err := logic.NewChatlogic().TestLoc(ctx, req)
	//response.Response(c, resp, err)
	//return
}

// GetFriendList
//
//	@Description: 建立Websocket连接，与好友聊天
//	@param c
func WebSocketHandler(c *gin.Context) {
	//ctx := zlog.GetCtxFromGin(c)
	//req, err := types.BindReq[types.Tes](c)
	//if err != nil {
	//	zlog.CtxErrorf(ctx, "WebSocketHandler request error: %v", err)
	//	response.NewResponse(c).Error(response.PARAM_NOT_VALID)
	//	return
	//}
	//zlog.CtxInfof(ctx, "WebSocketHandler request: %v", req)
	//resp, err := logic.NewChatlogic().TestLoc(ctx, req)
	//response.Response(c, resp, err)
	//return
}
