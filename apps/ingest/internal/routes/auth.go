package routes

import (
"database/sql"
"encoding/json"
"fmt"
"net/http"

"sentinelai/ingest/internal/auth"
"sentinelai/ingest/internal/store"
)

type tokenResp struct {
AccessToken string `json:"access_token"`
TokenType   string `json:"token_type"`
JTI         string `json:"jti"`
}

func Login(db *sql.DB) http.HandlerFunc {
return func(w http.ResponseWriter, r *http.Request) {
if err := r.ParseForm(); err != nil { http.Error(w, "bad form", 400); return }
email := r.FormValue("username")
pass := r.FormValue("password")

u, err := store.GetUserByEmail(r.Context(), db, email)
if err != nil || !u.IsActive || !store.VerifyPassword(pass, u.PasswordHash) {
http.Error(w, "bad credentials", http.StatusUnauthorized); return
}
scopes, err := store.CollectScopes(r.Context(), db, u.ID)
if err != nil { http.Error(w, "no scopes", http.StatusForbidden); return }

tok, jti, err := auth.Mint(fmt.Sprint(u.ID), scopes)
if err != nil { http.Error(w, "token error", 500); return }

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(tokenResp{AccessToken: tok, TokenType: "bearer", JTI: jti})
}
}
