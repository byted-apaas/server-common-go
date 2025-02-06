package http

import (
	"container/list"
	"sync"
	"time"
)

var (
	limiter = &RateLimiter{
		windowSize: time.Second, // 固定为 qps 限流
		maxRequest: -1,          // 默认不限流，在 AllowRequest 方法中实现
		requests:   list.New(),
		mutex:      sync.Mutex{},
	}
)

// RateLimiter 滑动窗口限流结构体
type RateLimiter struct {
	windowSize time.Duration
	maxRequest int
	requests   *list.List
	mutex      sync.Mutex
}

func (l *RateLimiter) ResetRateLimiter(maxRequest int) bool {
	if l.maxRequest == maxRequest {
		return false
	}

	l.maxRequest = maxRequest
	return true
}

// AllowRequest 判断是否允许请求
func (l *RateLimiter) AllowRequest() bool {
	if l.maxRequest <= 0 { // 不限流
		return true
	}

	l.mutex.Lock() // 加锁以保证并发安全
	defer l.mutex.Unlock()

	now := time.Now()
	// 移除滑动窗口外的请求
	for l.requests.Len() > 0 {
		front := l.requests.Front()
		if now.Sub(front.Value.(time.Time)) > l.windowSize {
			l.requests.Remove(front)
		} else {
			break
		}
	}

	// 如果请求数小于最大请求数，允许请求
	if l.requests.Len() < l.maxRequest {
		l.requests.PushBack(now)
		return true
	}

	return false
}
