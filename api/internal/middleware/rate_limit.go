package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/this-is-h/tikuAdapter-vercel/api/internal/registry/manager"
	"net/http"
)

// GlobalAPIRateLimit 全局API限流
func GlobalAPIRateLimit(c *gin.Context) {
	limiter := manager.GetManager().GetIPLimiter()
	congig := manager.GetManager().GetConfig()
	if congig.Limit.Enable && !limiter.GetLimiter(c.RemoteIP()).Allow() {
		c.AbortWithStatus(http.StatusTooManyRequests)
	}
}
