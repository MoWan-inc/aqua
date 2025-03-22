package util

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

// create a map to hold the rate limiters for each visitor and a mutex
var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

var tokenLimit = flag.Float64("token-limit", 1, "http serve token rate limit")
var tokenBursts = flag.Int("token-burst", 5, "http serve token burst")

func getVisitorLimiter(token string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, ok := visitors[token]
	if !ok {
		limiter = rate.NewLimiter(rate.Limit(*tokenLimit), *tokenBursts)
		visitors[token] = limiter
	}
	return limiter
}

func TokenLimit(next GinServerHandler) GinServerHandler {
	return func(ctx *gin.Context) (any, error) {
		// 从请求中获取token，如果没有带token则按照默认限流处理
		token, _ := getTokenFromCtx(ctx) //nolint:errcheck

		limiter := getVisitorLimiter(token.Token)
		if !limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests,
				gin.H{"msg": fmt.Sprintf("too many requests for token %s", token.Token)})
			return nil, nil
		}
		return next(ctx)
	}
}
