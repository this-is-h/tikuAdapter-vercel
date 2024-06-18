package registry

import (
	"github.com/this-is-h/tikuAdapter-vercel/api/configs"
	"github.com/this-is-h/tikuAdapter-vercel/api/pkg/ratelimit"
	"golang.org/x/time/rate"
)

// Limit get ratelimit instance
func Limit(cfg configs.Config) *ratelimit.IPRateLimiter {
	limit := cfg.Limit
	r := rate.Limit(float64(limit.Duration) / float64(limit.Requests))
	return ratelimit.NewIPRateLimiter(r, 1)
}
