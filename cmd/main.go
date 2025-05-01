package main

import (
	"TestCloudCamp/internal/ratelimit"
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"TestCloudCamp/internal/balancer"
	"TestCloudCamp/internal/config"
	"TestCloudCamp/internal/server"
)

func main() {
	// Загрузка конфига
	cfg, err := config.LoadConfig("./config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация бэкендов и HealthChecks
	var backends []*balancer.Backend
	for _, addr := range cfg.Backends {
		u, err := balancer.NewBackend(addr)
		if err != nil {
			log.Fatalf("Invalid backend URL %s: %v,", addr, err)
		}
		backends = append(backends, u)
	}

	pool := balancer.NewServerPool(backends)
	go func() {
		log.Println("Starting health checks...")
		pool.StartHealthChecks(5 * time.Second)
	}()

	// Инициализируем rate limiter
	r := rate.Limit(cfg.RateLimit.RPS)
	b := cfg.RateLimit.Burst
	rm := ratelimit.NewIPRateLimiter(r, b)

	// Создаём ReverseProxy и оборачиваем в middleware
	proxy := server.NewReverseProxy(pool)
	handler := rm.Middleware(proxy)

	// Graceful Shutdown
	ser := &http.Server{
		Addr:    cfg.Port,
		Handler: handler,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		<-ctx.Done()
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := ser.Shutdown(ctxTimeout); err != nil {
			log.Fatalf("Failed to shutdown server: %v", err)
		}
	}()

	fmt.Printf("Server started on port %v\n", cfg.Port)
	log.Fatal(ser.ListenAndServe())
}
