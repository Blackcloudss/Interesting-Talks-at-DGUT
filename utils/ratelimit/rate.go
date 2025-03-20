package ratelimit

import "golang.org/x/time/rate"

// @Title        ratelimit.go
// @Description
// @Create       XdpCs 2025-03-20 下午7:58
// @Update       XdpCs 2025-03-20 下午7:58

// rate基于令牌桶算法实现流量控制，防止系统过载
// rate.Limit表示每秒允许的请求数
type RateLimiter struct {
	limiter *rate.Limiter
}

// NewRateLimiter() 创建一个限流器，参数r表示每秒允许的请求数，b表示桶的容量
// 初始化示例：每秒10个令牌，突发容量20个
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(r, b),
	}
}

// Allow() 返回布尔值，表示当前是否允许请求通过（消耗1个令牌）
func (rl *RateLimiter) Allow() bool {
	return rl.limiter.Allow()
}
