package auth

import (
"os"
"time"

"github.com/golang-jwt/jwt/v5"
"github.com/google/uuid"
)

type Claims struct {
Sub   string   `json:"sub"`
Scope []string `json:"scope"`
jwt.RegisteredClaims
}

func ttl() time.Duration {
if v := os.Getenv("ACCESS_TTL_MIN"); v != "" {
if d, err := time.ParseDuration(v + "m"); err == nil { return d }
}
return 30 * time.Minute
}

func secret() []byte { return []byte(os.Getenv("JWT_SECRET")) }

func Mint(sub string, scopes []string) (string, string, error) {
jti := uuid.NewString()
now := time.Now()
claims := Claims{
Sub:   sub,
Scope: scopes,
RegisteredClaims: jwt.RegisteredClaims{
IssuedAt:  jwt.NewNumericDate(now),
ExpiresAt: jwt.NewNumericDate(now.Add(ttl())),
ID:        jti,
},
}
t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
signed, err := t.SignedString(secret())
return signed, jti, err
}

func Parse(raw string) (*Claims, error) {
tok, err := jwt.ParseWithClaims(raw, &Claims{}, func(t *jwt.Token) (interface{}, error) {
return secret(), nil
})
if err != nil { return nil, err }
if c, ok := tok.Claims.(*Claims); ok && tok.Valid { return c, nil }
return nil, jwt.ErrTokenInvalidClaims
}
