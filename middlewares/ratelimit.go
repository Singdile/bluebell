package middlewares

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// IPRateLimiter 按IP分组的限流器
type IPRateLimiter struct {
	ips sync.Map //map[string]*ratelimit.Bucket
}

// NewIPRateLimiter 创建IP限流器
func NewIPRateLimiter() *IPRateLimiter {
	return &IPRateLimiter{}
}

// 获取指定IP的令牌桶
// rate: 每秒放入的令牌数目
// capacity: 桶容量
func (i *IPRateLimiter) GetBucket(ip string, rate float64, capacity int64) *ratelimit.Bucket {
	// 首先查找缓存是否有对应的IP的令牌桶
	if bucket, ok := i.ips.Load(ip); ok {
		return bucket.(*ratelimit.Bucket)
	}

	// 创建IP对应的令牌桶
	bucket := ratelimit.NewBucketWithRate(rate, capacity)
	actual, _ := i.ips.LoadOrStore(ip, bucket)

	return actual.(*ratelimit.Bucket)
}

// RateLimitMiddleware 限流中间件
// rate: 每秒允许的请求数目
// capacity:突发容量
func RateLimitMiddleware(rate float64, capacity int64) gin.HandlerFunc {
	limiter := NewIPRateLimiter()

	return func(ctx *gin.Context) {
		//获取用户ip
		ip := ctx.ClientIP()

		//获取该ip对应的令牌桶
		bucket := limiter.GetBucket(ip, rate, capacity)

		//尝试获取一个令牌
		taken := bucket.TakeAvailable(1)

		//无令牌立即返回
		if taken == 0 {
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "请求过于频繁，请稍后再试",
			})
			ctx.Abort()
			return
		}

		//有令牌，继续请求
		ctx.Next()

	}
}
