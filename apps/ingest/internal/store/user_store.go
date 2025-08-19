package store

import (
"context"
"database/sql"
"errors"

"golang.org/x/crypto/bcrypt"
)

type User struct {
ID           int64
Email        string
PasswordHash string
IsActive     bool
}

func GetUserByEmail(ctx context.Context, db *sql.DB, email string) (*User, error) {
u := &User{}
err := db.QueryRowContext(ctx,
`SELECT id, email, password_hash, is_active FROM users WHERE email=$1`, email).
Scan(&u.ID, &u.Email, &u.PasswordHash, &u.IsActive)
if err != nil { return nil, err }
return u, nil
}

func VerifyPassword(password, hash string) bool {
return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func CollectScopes(ctx context.Context, db *sql.DB, userID int64) ([]string, error) {
rows, err := db.QueryContext(ctx, `
SELECT DISTINCT jsonb_array_elements_text(r.scopes->'scopes')
FROM roles r JOIN user_roles ur ON ur.role_id = r.id
WHERE ur.user_id = $1`, userID)
if err != nil { return nil, err }
defer rows.Close()

var scopes []string
for rows.Next() {
var s string
if err := rows.Scan(&s); err != nil { return nil, err }
scopes = append(scopes, s)
}
if len(scopes) == 0 { return nil, errors.New("no scopes") }
return scopes, nil
}
