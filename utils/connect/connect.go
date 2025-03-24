package connect

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// @Title        connect.go
// @Description
// @Create       XdpCs 2025-03-20 下午7:50
// @Update       XdpCs 2025-03-20 下午7:50

/*
WebSocket连接管理模块
主要功能：
	维护活跃连接集合
	提供线程安全的连接操作
	心跳保活机制
	连接数统计
*/

// ConnectionManager 连接管理器
// 解决：多协程并发访问安全问题、连接生命周期管理
type ConnectionManager struct {
	Clients map[int64]*websocket.Conn // 用户ID到连接的映射
	Mutex   sync.Mutex                // 保证并发安全的互斥锁
}

// NewConnectionManager 创建连接管理器
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		Clients: make(map[int64]*websocket.Conn),
	}
}

// AddClient 添加连接
func (cm *ConnectionManager) AddClient(userID int64, conn *websocket.Conn) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	cm.Clients[userID] = conn
}

// RemoveClient 删除连接
func (cm *ConnectionManager) RemoveClient(id int64) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	delete(cm.Clients, id)
}

// CheckConnections 心跳检测协程 -- 服务端主动检测
// 解决：网络闪断检测、僵尸连接清理
// 每隔15秒发送心跳包，如果连接超时，则删除连接并关闭连接
func (cm *ConnectionManager) CheckConnections() {
	go func() {
		// 创建定时器，定期向通道C发送时间值
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for {
			//  等待定时器超时，从ticker的通道C中接收一个值。
			//	由于通道操作在没有数据时会阻塞，所以这里会等待直到下一次定时触发，也就是等待一个时间间隔
			<-ticker.C
			cm.Mutex.Lock()
			// 遍历连接集合，发送心跳包
			for id, conn := range cm.Clients {
				// 设置写超时 ：防止心跳包发送阻塞（网络故障时快速失败
				conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
				// 发送心跳包: 首先尝试发送Ping消息，然后检查是否有错误（超时失败）。如果有错误，就执行移除客户端和关闭连接；否则流程结束。
				if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second)); err != nil {
					cm.RemoveClient(id)
					conn.Close()
				}
			}
			cm.Mutex.Unlock()
		}
	}()
}

// ConnectionCount 连接数统计
// 解决：系统负载监控、连接限制
func (cm *ConnectionManager) ConnectionCount() int {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	return len(cm.Clients)
}
