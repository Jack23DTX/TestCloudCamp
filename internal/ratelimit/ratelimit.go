package ratelimit

import (
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

// IPRateLimiter хранит лимитеры для каждого IP
type IPRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	r        rate.Limit
	b        int
}

// NewIPRateLimiter создает новый лимитер
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		b:        b,
	}
}

// getLimiter возвращает существующий или создает новый лимитер для IP
func (i *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	lim, ok := i.limiters[ip]
	i.mu.RUnlock()
	if ok {
		return lim
	}

	i.mu.Lock()
	defer i.mu.Unlock()
	// двойная проверка в случае гонки
	if lim, ok = i.limiters[ip]; ok {
		return lim
	}

	limiter := rate.NewLimiter(i.r, i.b)
	i.limiters[ip] = limiter

	return limiter
}

// Middleware проверяет лимит для каждого запроса
func (i *IPRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем IP клиента (учитываем X-Forwarded-For)
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.RemoteAddr
		}

		if !i.getLimiter(ip).Allow() {
			http.Error(w, `{"code": 429, "message": "Rate limit exceeded"}`, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
