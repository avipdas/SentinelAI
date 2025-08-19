package auth

import (
"context"
"net/http"
"strings"
)

type ctxKey int
const claimsKey ctxKey = 1

func WithAuth(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
h := r.Header.Get("Authorization")
if !strings.HasPrefix(strings.ToLower(h), "bearer ") {
http.Error(w, "missing bearer", http.StatusUnauthorized); return
}
raw := strings.TrimSpace(h[len("Bearer "):])
c, err := Parse(raw)
if err != nil {
http.Error(w, "invalid token", http.StatusUnauthorized); return
}
ctx := context.WithValue(r.Context(), claimsKey, c)
next.ServeHTTP(w, r.WithContext(ctx))
})
}

func RequireScopes(scopes ...string) func(http.Handler) http.Handler {
return func(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
val := r.Context().Value(claimsKey)
c, ok := val.(*Claims)
if !ok { http.Error(w, "unauth", http.StatusUnauthorized); return }
have := map[string]struct{}{}
for _, s := range c.Scope { have[s] = struct{}{} }
for _, need := range scopes {
if _, ok := have[need]; !ok {
http.Error(w, "forbidden", http.StatusForbidden); return
}
}
next.ServeHTTP(w, r)
})
}
}
