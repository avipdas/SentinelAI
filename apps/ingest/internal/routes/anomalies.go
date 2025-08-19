package routes

import (
"encoding/json"
"net/http"

"sentinelai/ingest/internal/auth"
)

func createAnomaly() http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// TODO: replace stub with your real insert into anomalies
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string]any{"ok": true})
})
}

func CreateAnomalyProtected() http.Handler {
return auth.WithAuth(auth.RequireScopes("ingest")(createAnomaly()))
}
