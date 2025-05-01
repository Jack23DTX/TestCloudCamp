package main

import (
	"context"
	"fmt"
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
	cfg, err := config.LoadConfig("./config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// тут будет запуск серверов из бэкендов (если они локальные)
	//lb := loadBalancer

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ser := &http.Server{
		Addr:    cfg.Port,
		Handler: server.NewReverseProxy(pool),
	}

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
