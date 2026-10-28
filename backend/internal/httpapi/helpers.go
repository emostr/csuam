package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"csuam/backend/internal/db"
)

type ctxKey string

const userKey ctxKey = "user"

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func urlID(r *http.Request, param string) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, param), 10, 64)
}

func currentUser(ctx context.Context) db.User {
	u, _ := ctx.Value(userKey).(db.User)
	return u
}

func isNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func isModerator(u db.User) bool {
	return u.Role == db.RoleLibrarian || u.Role == db.RoleHeadTeacher
}
