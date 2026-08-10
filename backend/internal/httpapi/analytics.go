package httpapi

import (
	"net/http"

	"csuam/backend/internal/db"
)

func (s *Server) handleAnalyticsSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	total, err := s.q.CountMaterials(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось посчитать материалы")
		return
	}
	byCategory, err := s.q.CountMaterialsByCategory(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось посчитать категории")
		return
	}
	monthly, err := s.q.MonthlyAdditions(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось получить динамику")
		return
	}
	pending, _ := s.q.CountPendingMaterials(ctx)
	activeLoans, _ := s.q.CountActiveLoans(ctx)
	overdueLoans, _ := s.q.CountOverdueLoans(ctx)

	if byCategory == nil {
		byCategory = []db.CategoryCount{}
	}
	if monthly == nil {
		monthly = []db.MonthCount{}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"total":         total,
		"by_category":   byCategory,
		"monthly":       monthly,
		"pending":       pending,
		"active_loans":  activeLoans,
		"overdue_loans": overdueLoans,
	})
}
