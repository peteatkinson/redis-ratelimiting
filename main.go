package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/rwxpeter/sliding-rate-limit/limiters"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		w.WriteHeader(http.StatusOK)
	})

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	rateLimiter := limiters.New(client, time.Minute*10, 5)

	router.Use(limiters.HttpRateLimiter(rateLimiter, limiters.RealIP()))

	server := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8002",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
