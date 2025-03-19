package api

import (
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/jwt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

// @Description: 加载历史消息
// @param c
func GetHistoryMessage(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	UserID := jwt.GetUserId(c)
	req, err := types.BindReq[types.GetMessageReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetMessageHistory request error: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "GetMessageHistory request: %v", req)
	resp, err := logic.NewChatlogic().GetMessagesHistory(ctx, UserID, req)
	response.Response(c, resp, err)
	return
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients = make(map[int64]*websocket.Conn)
	mutex   sync.Mutex
)

// GetFriendList
//
//	@Description: 建立Websocket连接，与好友聊天
//	@param c
func WebSocketHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	UserId := jwt.GetUserId(c)
	// 升级 HTTP 连接为 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zlog.CtxErrorf(ctx, "WebSocketHandler error: %v", err)
		response.NewResponse(c).Error(response.WEBSOCKET_UPGRADE_FAIL)
		return
	}
	defer conn.Close()

	// 注册连接
	mutex.Lock()
	clients[UserId] = conn
	mutex.Unlock()

	// 消息处理循环
	for {
		var msg types.WSMessage
		// 读取客户端消息
		if err = conn.ReadJSON(&msg); err != nil {
			zlog.CtxErrorf(ctx, "WebSocketHandler error: %v", err)
			response.NewResponse(c).Error(response.WEBSOCKET_READ_COMMENT_FAIL)
			return
		}
		// 检测接收者ID
		if msg.To == 0 {
			response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		}
		// 发送消息给接收者
		if err = logic.NewChatlogic().SendMessage(ctx, UserId, msg); err != nil {
			zlog.CtxErrorf(ctx, "WebSocketHandler error: %v", err)
			err = conn.WriteJSON(types.WSError{ERROR: err.Error()})
			if err != nil {
				zlog.CtxErrorf(ctx, "WebSocketHandler error: %v", err)
				response.NewResponse(c).Error(response.WEBSOCKET_WRITE_COMMENT_FAIL)
				return
			}
			return
		}
	}
}
