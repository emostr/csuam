package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"csuam/backend/internal/auth"
)

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(auth.CookieName)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "требуется вход")
			return
		}
		uid, err := auth.ParseToken(s.cfg.JWTSecret, cookie.Value)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "недействительная сессия")
			return
		}
		user, err := s.q.GetUserByID(r.Context(), uid)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "пользователь не найден")
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userKey, user)))
	})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "неверный формат запроса")
		return
	}
	user, err := s.q.GetUserByUsername(r.Context(), req.Username)
	if err != nil || !auth.CheckPassword(user.PasswordHash, req.Password) {
		writeError(w, http.StatusUnauthorized, "неверный логин или пароль")
		return
	}
	token, err := auth.IssueToken(s.cfg.JWTSecret, user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось создать сессию")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     auth.CookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
	writeJSON(w, http.StatusOK, user)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     auth.CookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, currentUser(r.Context()))
}

func (s *Server) handleChangePassword(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r.Context())
	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "неверный формат запроса")
		return
	}
	if !auth.CheckPassword(user.PasswordHash, req.CurrentPassword) {
		writeError(w, http.StatusForbidden, "текущий пароль указан неверно")
		return
	}
	if len(req.NewPassword) < 6 {
		writeError(w, http.StatusBadRequest, "новый пароль должен быть не короче 6 символов")
		return
	}
	hash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось обработать пароль")
		return
	}
	if err := s.q.UpdateUserPassword(r.Context(), user.ID, hash); err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось сменить пароль")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.q.ListUsers(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось получить пользователей")
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func (s *Server) handleNotifications(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r.Context())
	resp := map[string]int64{"overdue_loans": 0, "pending_materials": 0}
	if isModerator(user) {
		if n, err := s.q.CountOverdueLoans(r.Context()); err == nil {
			resp["overdue_loans"] = n
		}
		if n, err := s.q.CountPendingMaterials(r.Context()); err == nil {
			resp["pending_materials"] = n
		}
	}
	writeJSON(w, http.StatusOK, resp)
}
