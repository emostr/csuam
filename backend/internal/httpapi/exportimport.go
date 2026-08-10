package httpapi

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"csuam/backend/internal/db"
)

type materialCard struct {
	XMLName       xml.Name `json:"-" xml:"material"`
	Title         string   `json:"title" xml:"title"`
	Description   string   `json:"description" xml:"description"`
	Category      string   `json:"category" xml:"category"`
	Condition     string   `json:"condition" xml:"condition"`
	Location      string   `json:"location" xml:"location"`
	OriginDate    string   `json:"origin_date,omitempty" xml:"origin_date,omitempty"`
	Content       string   `json:"content,omitempty" xml:"content,omitempty"`
	ContentFormat string   `json:"content_format,omitempty" xml:"content_format,omitempty"`
	FileName      string   `json:"file_name,omitempty" xml:"file_name,omitempty"`
	CreatedAt     string   `json:"created_at" xml:"created_at"`
}

func cardFromMaterial(m db.Material) materialCard {
	card := materialCard{
		Title:       m.Title,
		Description: m.Description,
		Category:    m.Category,
		Condition:   m.Condition,
		Location:    m.Location,
		CreatedAt:   m.CreatedAt.Format("2006-01-02"),
	}
	if m.OriginDate != nil {
		card.OriginDate = m.OriginDate.Format("2006-01-02")
	}
	if m.Content != nil {
		card.Content = *m.Content
	}
	if m.ContentFormat != nil {
		card.ContentFormat = *m.ContentFormat
	}
	if m.FileName != nil {
		card.FileName = *m.FileName
	}
	return card
}

func (s *Server) handleExportMaterial(w http.ResponseWriter, r *http.Request) {
	m, ok := s.getAccessibleMaterial(w, r)
	if !ok {
		return
	}
	card := cardFromMaterial(m)
	format := r.URL.Query().Get("format")

	if format == "xml" {
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=material-%d.xml", m.ID))
		_, _ = w.Write([]byte(xml.Header))
		enc := xml.NewEncoder(w)
		enc.Indent("", "  ")
		_ = enc.Encode(card)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=material-%d.json", m.ID))
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(card)
}

func (s *Server) handleImportMaterial(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r.Context())

	var body []byte
	ct := r.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "multipart/form-data") {
		if err := r.ParseMultipartForm(16 << 20); err != nil {
			writeError(w, http.StatusBadRequest, "не удалось прочитать форму")
			return
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			writeError(w, http.StatusBadRequest, "приложите файл карточки (JSON или XML)")
			return
		}
		defer file.Close()
		body, err = io.ReadAll(io.LimitReader(file, 16<<20))
		if err != nil {
			writeError(w, http.StatusBadRequest, "не удалось прочитать файл")
			return
		}
	} else {
		var err error
		body, err = io.ReadAll(io.LimitReader(r.Body, 16<<20))
		if err != nil {
			writeError(w, http.StatusBadRequest, "не удалось прочитать запрос")
			return
		}
	}

	var card materialCard
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		writeError(w, http.StatusBadRequest, "файл пуст")
		return
	}
	if trimmed[0] == '{' {
		if err := json.Unmarshal(trimmed, &card); err != nil {
			writeError(w, http.StatusBadRequest, "не удалось разобрать JSON")
			return
		}
	} else {
		if err := xml.Unmarshal(trimmed, &card); err != nil {
			writeError(w, http.StatusBadRequest, "не удалось разобрать XML")
			return
		}
	}

	if strings.TrimSpace(card.Title) == "" {
		writeError(w, http.StatusBadRequest, "в карточке не указано название")
		return
	}
	if !validCategories[card.Category] {
		writeError(w, http.StatusBadRequest, "в карточке указана неверная категория")
		return
	}
	if user.Role == db.RoleLibrarian && card.Category == db.CategoryDocuments {
		writeError(w, http.StatusForbidden, "библиотекарь не может добавлять документы")
		return
	}

	status := db.StatusApproved
	if user.Role == db.RoleTeacher {
		status = db.StatusPending
	}

	p := db.CreateMaterialParams{
		Title:         strings.TrimSpace(card.Title),
		Description:   card.Description,
		Category:      card.Category,
		Status:        status,
		Condition:     card.Condition,
		Location:      card.Location,
		OriginDate:    parseDate(card.OriginDate),
		Content:       optStr(card.Content),
		ContentFormat: optStr(card.ContentFormat),
		CreatedBy:     &user.ID,
	}
	id, err := s.q.CreateMaterial(r.Context(), p)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось импортировать материал")
		return
	}
	m, err := s.q.GetMaterial(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "материал создан, но не удалось его прочитать")
		return
	}
	writeJSON(w, http.StatusCreated, m)
}
