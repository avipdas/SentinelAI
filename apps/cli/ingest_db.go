package main

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    "github.com/go-redis/redis/v8"
    _ "github.com/lib/pq"
)

// use ctx from main.go

func StartDBWriter() {
    // Redis connection
    rdb := redis.NewClient(&redis.Options{
        Addr: "redis:6379",
    })

    // Postgres connection
    connStr := "postgres://sentinel:sentinel123@db:5432/sentinelai?sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Postgres connection error:", err)
    }
    defer db.Close()

    log.Println("DB writer running. Moving logs from Redis to Postgres...")

    for {
        msg, err := rdb.LPop(ctx, "logs").Result()
        if err == redis.Nil {
            time.Sleep(2 * time.Second) // wait if no logs
            continue
        } else if err != nil {
            log.Println("Redis read error:", err)
            continue
        }

        _, err = db.Exec("INSERT INTO logs (message) VALUES ($1)", msg)
        if err != nil {
            log.Println("Postgres insert error:", err)
        } else {
            fmt.Println("Inserted log into Postgres:", msg)
        }
    }
}
