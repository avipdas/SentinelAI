package store

import (
"database/sql"
"os"

_ "github.com/lib/pq"
)

func MustOpen() *sql.DB {
dsn := os.Getenv("DATABASE_URL")
if dsn == "" {
dsn = "postgres://sentinel:sentinel@localhost:5432/sentinelai?sslmode=disable"
}
db, err := sql.Open("postgres", dsn)
if err != nil { panic(err) }
if err := db.Ping(); err != nil { panic(err) }
return db
}
