package main

import (
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/go-redis/redis/v8"
    "context"
)

var ctx = context.Background()

func main() {
    rdb := redis.NewClient(&redis.Options{
        Addr: "redis:6379",
    })

    http.HandleFunc("/ingest", func(w http.ResponseWriter, r *http.Request) {
        msg := r.URL.Query().Get("msg")
        if msg == "" {
            http.Error(w, "Missing msg param", http.StatusBadRequest)
            return
        }

        err := rdb.LPush(ctx, "logs", msg).Err()
        if err != nil {
            http.Error(w, "Failed to push to Redis", http.StatusInternalServerError)
            return
        }

        fmt.Fprintf(w, "Log received at %s: %s", time.Now().Format(time.RFC3339), msg)
    })

	go StartDBWriter() // Start the Redis → Postgres worker

    log.Println("Ingest service running on :8080")
    http.ListenAndServe(":8080", nil)
}
