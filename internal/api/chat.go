package api

import (
	"bytes"
	"encoding/json"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/connect"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/ratelimit"

	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/logic"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/response"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/internal/types"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/utils/jwt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/klauspost/compress/gzip"
)

const (
	WEBSOCKET_UPGRADE_FAIL       = "60007, websocket升级失败"
	WEBSOCKET_READ_COMMENT_FAIL  = "60008, websocket读取消息失败"
	WEBSOCKET_WRITE_COMMENT_FAIL = "60009, websocket发送消息失败"
)

var (
	upgrader = websocket.Upgrader{
		EnableCompression: true, // 启用压缩支持
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// 全局速率限制（每秒 100 个消息） 流量控制
	Limiter = ratelimit.NewRateLimiter(100, 200)

	// 注册一个连接管理器
	CM = connect.NewConnectionManager()

	// 添加消息确认处理
	ackChannels = sync.Map{}
)

/*
WebSocket API接口模块
主要功能：
1. HTTP升级WebSocket协议
2. 连接速率限制
3. 消息编解码处理
4. 离线消息同步
5. 异常处理
*/

// @Description: 加载历史消息
// @param c
func GetHistoryMessage(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	UserID := jwt.GetUserId(c)
	req, err := types.BindReq[types.GetMessageReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetMessageHistory request error: %v", err)
		return
	}
	zlog.CtxInfof(ctx, "GetMessageHistory request: %v", req)
	resp, err := logic.NewChatlogic().GetMessagesHistory(ctx, UserID, req)
	response.Response(c, resp, err)
	return
}

// WebSocketHandler 连接入口
// 解决：协议升级、流量控制、资源清理
func WebSocketHandler(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	UserId := jwt.GetUserId(c)
	// 升级 HTTP 连接为 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zlog.CtxErrorf(ctx, "WebSocketHandler error: %v", err)
		_ = conn.WriteJSON(types.WSError{ERROR: WEBSOCKET_UPGRADE_FAIL + err.Error()})
		return
	}

	// 令牌桶算法
	if !Limiter.Allow() {
		_ = conn.WriteJSON(types.WSError{ERROR: "请求过于频繁"})
		conn.Close()
		return
	}

	// 限制连接数（修改为使用统计方法）
	if CM.ConnectionCount() > 1000 {
		_ = conn.WriteJSON(types.WSError{ERROR: "系统繁忙"})
		conn.Close()
		return
	}

	defer func() {
		// 删除连接
		CM.RemoveClient(UserId)
		conn.Close()
	}()

	// 替换原有注册连接
	CM.Mutex.Lock()
	if oldConn, exists := CM.Clients[UserId]; exists {
		oldConn.Close()
		delete(CM.Clients, UserId)
	}
	CM.Mutex.Unlock()
	// 注册新连接
	CM.AddClient(UserId, conn)

	// 获取离线消息
	if msgs, err := logic.NewChatlogic().GetOfflineMessages(UserId); err == nil {
		for _, msg := range msgs {
			msg.Type = global.MESSAGE // 确保类型正确
			_ = conn.WriteJSON(msg)
		}
	}

	// 心跳检测
	// 设置 Pong 响应处理器 -- 接受客户端 Pong 响应处理
	conn.SetPongHandler(func(string) error {
		// 收到 Pong 包时，重置读超时时间窗口，自动延长60秒读超时
		err = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		if err != nil {
			zlog.CtxErrorf(ctx, "WebSocket设置读超时失败: %v", err)
		}
		return nil
	})

	// 添加初始读超时设置
	//1.防止僵尸连接占用资源（60秒内无任何数据读取自动断开）；
	//2.配合心跳机制实现活性检测（收到Pong响应后重置倒计时，维持长连接）
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	// 启动独立协程发送心跳包
	CM.CheckConnections()

	// 消息处理循环关键部分：
	// 解决：消息超时控制、压缩处理、确认机制
	for {
		// 设置双超时机制：写超时10s，读超时60s
		if err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
			zlog.CtxErrorf(ctx, "设置写超时失败: %v", err)
			return
		}

		// 读取并解压消息
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			zlog.CtxErrorf(ctx, "WebSocket读取消息失败: %v", err)
			_ = conn.WriteJSON(types.WSError{ERROR: WEBSOCKET_READ_COMMENT_FAIL + err.Error()})
			return
		}

		// 自动解压gzip消息
		if messageType == websocket.BinaryMessage {
			gr, err := gzip.NewReader(bytes.NewReader(p))
			if err != nil {
				zlog.CtxErrorf(ctx, "创建解压读取器失败: %v", err)
				continue
			}
			defer gr.Close()
			// 读取并解压消息
			decompressed, err := io.ReadAll(gr)
			if err != nil {
				zlog.CtxErrorf(ctx, "解压消息失败: %v", err)
				continue
			}
			// 更新消息内容
			p = decompressed
		}
		// 解析用户要发送的消息
		var msg types.WSMessageReq
		if err = json.Unmarshal(p, &msg); err != nil {
			zlog.CtxErrorf(ctx, "消息解析失败: %v", err)
			continue
		}
		// 消息确认处理， 分离普通消息与确认消息处理
		// 普通消息携带业务数据，解压处理； ACK仅包含控制信息， 触发确认
		//subgraph 消息处理循环
		//M[读取消息] --> N{消息类型}
		//N -->|普通消息| O[解压处理]
		//N -->|ACK消息| P[触发确认]
		//O --> Q[好友验证]
		//Q --> R[存储消息]
		//R --> S{接收方在线?}
		//S -->|在线| T[存储为delivered 在线消息，加入消息队列]
		//S -->|离线| U[存储为offline 离线消息]
		//T --> V[批量压缩发送]
		//V --> W[记录发送时间]
		//W --> X[等待用户发来的 ACK]
		//X -->|30秒未收到| Y[重发消息]
		//X -->|收到ACK| Z[更新状态]
		//end
		//
		//subgraph 确认机制
		//P --> AA{存在消息通道?}
		//AA -->|是| AB[关闭通道]
		//AB --> AC[更新为ACKNOWLEDGED]
		//AA -->|否| AD[记录异常]
		//AC --> AE[清理确认记录]
		//end
		if msg.Type == global.ACK {
			// 收到ACK消息时，会通过ackChannels.Load找到对应msgID的通道
			if ch, ok := ackChannels.Load(msg.MsgID); ok {
				// 通知消息已确认， 关闭通道，释放资源
				// 关闭通道会立即返回零值，此时会触发ackChan返回的通道的case <-ackChan(msgID)分支，完成消息确认
				close(ch.(chan struct{}))
			}
			continue
		}

		// 检测接收者ID
		if msg.To == 0 {
			response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		}

		// 发送消息给接收者
		if err = logic.NewChatlogic().SendMessage(ctx, UserId, msg, CM); err != nil {
			zlog.CtxErrorf(ctx, "WebSocket发送消息失败: %v", err)
			_ = conn.WriteJSON(types.WSError{ERROR: WEBSOCKET_WRITE_COMMENT_FAIL + err.Error()})
			return
		}

		// 发送响应时增加超时控制
		if err = conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
			zlog.CtxErrorf(ctx, "设置响应超时失败: %v", err)
			return
		}
	}
}
