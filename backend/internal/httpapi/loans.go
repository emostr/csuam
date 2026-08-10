package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"csuam/backend/internal/db"
)

func (s *Server) requireLibrarian(w http.ResponseWriter, r *http.Request) bool {
	user := currentUser(r.Context())
	if !isModerator(user) {
		writeError(w, http.StatusForbidden, "прокат доступен библиотекарю и завучу")
		return false
	}
	return true
}

func (s *Server) handleListLoans(w http.ResponseWriter, r *http.Request) {
	if !s.requireLibrarian(w, r) {
		return
	}
	active := r.URL.Query().Get("archive") != "1"
	loans, err := s.q.ListLoans(r.Context(), active)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось получить список проката")
		return
	}
	if loans == nil {
		loans = []db.Loan{}
	}
	writeJSON(w, http.StatusOK, loans)
}

func (s *Server) handleCreateLoan(w http.ResponseWriter, r *http.Request) {
	if !s.requireLibrarian(w, r) {
		return
	}
	var req struct {
		MaterialID   int64  `json:"material_id"`
		BorrowerName string `json:"borrower_name"`
		Note         string `json:"note"`
		DueDate      string `json:"due_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "неверный формат запроса")
		return
	}
	if req.BorrowerName == "" {
		writeError(w, http.StatusBadRequest, "укажите, кто взял предмет")
		return
	}
	due, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "укажите срок возврата")
		return
	}
	if _, err := s.q.GetMaterial(r.Context(), req.MaterialID); err != nil {
		writeError(w, http.StatusNotFound, "материал не найден")
		return
	}
	id, err := s.q.CreateLoan(r.Context(), req.MaterialID, req.BorrowerName, req.Note, due)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось создать запись проката")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func (s *Server) handleReturnLoan(w http.ResponseWriter, r *http.Request) {
	if !s.requireLibrarian(w, r) {
		return
	}
	id, err := urlID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "неверный идентификатор")
		return
	}
	if err := s.q.ReturnLoan(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось отметить возврат")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
