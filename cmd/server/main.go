package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lucas-sachet/full-cycle-desafio-rate-limiter/internal/config"
	"github.com/lucas-sachet/full-cycle-desafio-rate-limiter/internal/limiter"
	"github.com/lucas-sachet/full-cycle-desafio-rate-limiter/internal/middleware"
	"github.com/lucas-sachet/full-cycle-desafio-rate-limiter/internal/storage"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	redisStrategy := storage.NewRedisStrategy(
		cfg.RedisHost,
		cfg.RedisPort,
		cfg.RedisPassword,
		cfg.RedisDB,
	)

	rl := limiter.NewRateLimiter(
		redisStrategy,
		cfg.DefaultIPLimit,
		cfg.DefaultTokenLimit,
		cfg.BlockDurationSeconds,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Wrap the mux with the rate limiter middleware
	handler := middleware.RateLimiterMiddleware(rl)(mux)

	serverAddr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server on %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
