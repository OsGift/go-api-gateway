package middleware

import (
	"sync"

	"golang.org/x/time/rate"
	"github.com/valyala/fasthttp"
)

var (
	limiterMap = make(map[string]*rate.Limiter)
	mu         sync.Mutex
	rateLimit  = rate.Limit(10) // 10 requests per second
	burstLimit = 5
)

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := limiterMap[ip]; !exists {
		limiterMap[ip] = rate.NewLimiter(rateLimit, burstLimit)
	}

	return limiterMap[ip]
}

func RateLimitMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ip := ctx.RemoteIP().String()
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			ctx.SetStatusCode(fasthttp.StatusTooManyRequests)
			ctx.SetBody([]byte("Too many requests"))
			return
		}

		next(ctx)
	}
}
