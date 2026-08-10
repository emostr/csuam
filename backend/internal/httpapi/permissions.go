package httpapi

import (
	"encoding/json"
	"net/http"

	"csuam/backend/internal/db"
)

func (s *Server) handleListPermissions(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r.Context())
	if user.Role != db.RoleHeadTeacher {
		writeError(w, http.StatusForbidden, "права доступа выдаёт завуч")
		return
	}
	id, err := urlID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "неверный идентификатор")
		return
	}
	perms, err := s.q.ListPermissions(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось получить список прав")
		return
	}
	if perms == nil {
		perms = []db.Permission{}
	}
	writeJSON(w, http.StatusOK, perms)
}

func (s *Server) handleGrantPermission(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r.Context())
	if user.Role != db.RoleHeadTeacher {
		writeError(w, http.StatusForbidden, "права доступа выдаёт завуч")
		return
	}
	id, err := urlID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "неверный идентификатор")
		return
	}
	var req struct {
		UserID int64 `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == 0 {
		writeError(w, http.StatusBadRequest, "укажите пользователя")
		return
	}
	if _, err := s.q.GetMaterial(r.Context(), id); err != nil {
		writeError(w, http.StatusNotFound, "материал не найден")
		return
	}
	if err := s.q.GrantPermission(r.Context(), id, req.UserID, user.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось выдать доступ")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleRevokePermission(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r.Context())
	if user.Role != db.RoleHeadTeacher {
		writeError(w, http.StatusForbidden, "права доступа выдаёт завуч")
		return
	}
	id, err := urlID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "неверный идентификатор")
		return
	}
	userID, err := urlID(r, "userId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "неверный идентификатор пользователя")
		return
	}
	if err := s.q.RevokePermission(r.Context(), id, userID); err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось отозвать доступ")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
