package ratelimiter

import (
	"sync"
	"time"
)

type FixWindowRateLimiter struct {
	mu          sync.Mutex
	count       map[int64]int
	countReject bool
	finish      chan struct{}
}

func NewFixWindowRateLimiter() *FixWindowRateLimiter {
	limiter := &FixWindowRateLimiter{
		mu:          sync.Mutex{},
		count:       map[int64]int{},
		countReject: false,
	}
	ticker := time.NewTicker(time.Second)
	go limiter.cleanKeys(ticker)
	return limiter
}

func (l *FixWindowRateLimiter) Close() {
	l.finish <- struct{}{}
}

func (l *FixWindowRateLimiter) cleanKeys(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			timeNow := time.Now().Unix()
			l.mu.Lock()
			for k := range l.count {
				if k < timeNow-2 {
					delete(l.count, k)
				}
			}
			if timeNow%1000 == 0 {
				newMap := make(map[int64]int)
				for k, v := range l.count {
					newMap[k] = v
				}
				l.count = newMap
			}
			l.mu.Unlock()
		case <-l.finish:
			return
		}
	}
}

func (l *FixWindowRateLimiter) Pass(n int) bool {
	timeNow := time.Now().Unix()
	l.mu.Lock()
	val := l.count[timeNow]
	result := n > val
	if l.countReject || result == true {
		l.count[timeNow] = val + 1
	}
	return result
}

func (l *FixWindowRateLimiter) SetCountReject(countReject bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.countReject = countReject
}
