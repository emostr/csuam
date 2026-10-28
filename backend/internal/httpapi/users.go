package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"

	"csuam/backend/internal/auth"
	"csuam/backend/internal/db"
)

var validRoles = map[string]bool{
	db.RoleTeacher: true, db.RoleLibrarian: true, db.RoleHeadTeacher: true,
}

func (s *Server) requireHeadTeacher(w http.ResponseWriter, r *http.Request) bool {
	if currentUser(r.Context()).Role != db.RoleHeadTeacher {
		writeError(w, http.StatusForbidden, "управление пользователями доступно только завучу")
		return false
	}
	return true
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	if !s.requireHeadTeacher(w, r) {
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		FullName string `json:"full_name"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "неверный формат запроса")
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	req.FullName = strings.TrimSpace(req.FullName)
	if req.Username == "" || req.FullName == "" {
		writeError(w, http.StatusBadRequest, "укажите логин и ФИО")
		return
	}
	if len(req.Password) < 6 {
		writeError(w, http.StatusBadRequest, "пароль должен быть не короче 6 символов")
		return
	}
	if !validRoles[req.Role] {
		writeError(w, http.StatusBadRequest, "неверная роль")
		return
	}
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось обработать пароль")
		return
	}
	user, err := s.q.CreateUser(r.Context(), req.Username, hash, req.FullName, req.Role)
	if err != nil {
		if isUniqueViolation(err) {
			writeError(w, http.StatusConflict, "такой логин уже занят")
			return
		}
		writeError(w, http.StatusInternalServerError, "не удалось создать пользователя")
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if !s.requireHeadTeacher(w, r) {
		return
	}
	id, err := urlID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "неверный идентификатор")
		return
	}
	if id == currentUser(r.Context()).ID {
		writeError(w, http.StatusBadRequest, "нельзя удалить собственную учётную запись")
		return
	}
	if _, err := s.q.GetUserByID(r.Context(), id); err != nil {
		writeError(w, http.StatusNotFound, "пользователь не найден")
		return
	}
	if err := s.q.DeleteUser(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось удалить пользователя")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleResetPassword(w http.ResponseWriter, r *http.Request) {
	if !s.requireHeadTeacher(w, r) {
		return
	}
	id, err := urlID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "неверный идентификатор")
		return
	}
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "неверный формат запроса")
		return
	}
	if len(req.Password) < 6 {
		writeError(w, http.StatusBadRequest, "пароль должен быть не короче 6 символов")
		return
	}
	if _, err := s.q.GetUserByID(r.Context(), id); err != nil {
		writeError(w, http.StatusNotFound, "пользователь не найден")
		return
	}
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось обработать пароль")
		return
	}
	if err := s.q.UpdateUserPassword(r.Context(), id, hash); err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось сменить пароль")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
