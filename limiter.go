package ratelimiter

type RateLimiter interface {
	Pass(n int) bool
	SetCountReject(countReject bool)
	Close()
}
