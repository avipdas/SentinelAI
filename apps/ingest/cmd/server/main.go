package main

import (
"log"
"net/http"

"sentinelai/ingest/internal/routes"
"sentinelai/ingest/internal/store"
)

func main() {
db := store.MustOpen()
mux := http.NewServeMux()

// Auth (POST form: username, password)
mux.Handle("/auth/token", routes.Login(db))

// Anomalies (JWT + "ingest" scope)
mux.Handle("/anomalies", routes.CreateAnomalyProtected())

addr := ":8000"
log.Println("ingest listening on", addr)
log.Fatal(http.ListenAndServe(addr, mux))
}
