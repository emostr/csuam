package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	qrcode "github.com/skip2/go-qrcode"

	"csuam/backend/internal/db"
)

var allowedExtensions = map[string]bool{
	".pdf": true, ".doc": true, ".docx": true, ".ppt": true, ".pptx": true,
	".odt": true, ".odp": true, ".png": true, ".jpg": true, ".jpeg": true,
	".gif": true, ".svg": true, ".bmp": true, ".webp": true, ".mp4": true,
	".mp3": true, ".wav": true, ".wmv": true, ".webm": true, ".html": true,
	".md": true,
}

var validCategories = map[string]bool{
	db.CategoryAwards: true, db.CategoryPhotos: true, db.CategoryVideos: true,
	db.CategoryLibrary: true, db.CategoryDocuments: true,
}

const maxUploadSize = 512 << 20

func (s *Server) canView(r *http.Request, u db.User, m db.Material) bool {
	if u.Role == db.RoleHeadTeacher {
		return true
	}
	if m.CreatedBy != nil && *m.CreatedBy == u.ID {
		return true
	}
	if m.Category == db.CategoryDocuments {
		ok, err := s.q.HasPermission(r.Context(), m.ID, u.ID)
		if err != nil || !ok {
			return false
		}
	}
	if m.Status != db.StatusApproved && u.Role != db.RoleLibrarian {
		return false
	}
	return true
}

func (s *Server) canManage(r *http.Request, u db.User, m db.Material) bool {
	switch u.Role {
	case db.RoleHeadTeacher:
		return true
	case db.RoleLibrarian:
		if m.Category != db.CategoryDocuments {
			return true
		}
		ok, err := s.q.HasPermission(r.Context(), m.ID, u.ID)
		return err == nil && ok
	default:
		return m.CreatedBy != nil && *m.CreatedBy == u.ID
	}
}

func parseDate(v string) *time.Time {
	if v == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return nil
	}
	return &t
}

func optStr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func (s *Server) handleListMaterials(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r.Context())
	qp := r.URL.Query()

	p := db.ListMaterialsParams{
		Query:      qp.Get("q"),
		FromDate:   parseDate(qp.Get("from")),
		ToDate:     parseDate(qp.Get("to")),
		ViewerID:   user.ID,
		ViewerRole: user.Role,
	}
	if c := qp.Get("category"); c != "" && validCategories[c] {
		p.Category = &c
	}
	if st := qp.Get("status"); st == db.StatusPending || st == db.StatusApproved || st == db.StatusRejected {
		p.Status = &st
	}

	items, err := s.q.ListMaterials(r.Context(), p)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось получить материалы")
		return
	}
	if items == nil {
		items = []db.Material{}
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleCreateMaterial(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r.Context())

	status := db.StatusApproved
	if user.Role == db.RoleTeacher {
		status = db.StatusPending
	}

	p := db.CreateMaterialParams{
		Status:    status,
		CreatedBy: &user.ID,
	}

	var headers []*multipart.FileHeader

	ct := r.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "multipart/form-data") {
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			writeError(w, http.StatusBadRequest, "не удалось прочитать форму (слишком большой файл?)")
			return
		}
		p.Title = strings.TrimSpace(r.FormValue("title"))
		p.Description = r.FormValue("description")
		p.Category = r.FormValue("category")
		p.Condition = r.FormValue("condition")
		p.Location = r.FormValue("location")
		p.OriginDate = parseDate(r.FormValue("origin_date"))
		p.Content = optStr(r.FormValue("content"))
		p.ContentFormat = optStr(r.FormValue("content_format"))
		headers = formFileHeaders(r)
	} else {
		var req struct {
			Title         string `json:"title"`
			Description   string `json:"description"`
			Category      string `json:"category"`
			Condition     string `json:"condition"`
			Location      string `json:"location"`
			OriginDate    string `json:"origin_date"`
			Content       string `json:"content"`
			ContentFormat string `json:"content_format"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "неверный формат запроса")
			return
		}
		p.Title = strings.TrimSpace(req.Title)
		p.Description = req.Description
		p.Category = req.Category
		p.Condition = req.Condition
		p.Location = req.Location
		p.OriginDate = parseDate(req.OriginDate)
		p.Content = optStr(req.Content)
		p.ContentFormat = optStr(req.ContentFormat)
	}

	if p.Title == "" {
		writeError(w, http.StatusBadRequest, "укажите название материала")
		return
	}
	if !validCategories[p.Category] {
		writeError(w, http.StatusBadRequest, "неверная категория")
		return
	}
	if user.Role == db.RoleLibrarian && p.Category == db.CategoryDocuments {
		writeError(w, http.StatusForbidden, "библиотекарь не может добавлять документы")
		return
	}

	uploads, err := s.uploadFiles(r.Context(), headers)
	if err != nil {
		if errors.Is(err, errUnsupportedFormat) {
			writeError(w, http.StatusBadRequest, "формат файла не поддерживается")
		} else {
			writeError(w, http.StatusInternalServerError, "не удалось сохранить файл в хранилище")
		}
		return
	}

	id, err := s.q.CreateMaterial(r.Context(), p)
	if err != nil {
		s.dropUploads(r.Context(), uploads)
		writeError(w, http.StatusInternalServerError, "не удалось создать материал")
		return
	}
	if err := s.attachFiles(r.Context(), id, uploads); err != nil {
		s.dropUploads(r.Context(), uploads)
		writeError(w, http.StatusInternalServerError, "материал создан, но файлы не удалось прикрепить")
		return
	}
	m, err := s.q.GetMaterial(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "материал создан, но не удалось его прочитать")
		return
	}
	writeJSON(w, http.StatusCreated, m)
}

func (s *Server) getAccessibleMaterial(w http.ResponseWriter, r *http.Request) (db.Material, bool) {
	id, err := urlID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "неверный идентификатор")
		return db.Material{}, false
	}
	m, err := s.q.GetMaterial(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, "материал не найден")
		} else {
			writeError(w, http.StatusInternalServerError, "не удалось получить материал")
		}
		return db.Material{}, false
	}
	if !s.canView(r, currentUser(r.Context()), m) {
		writeError(w, http.StatusForbidden, "нет доступа к этому материалу")
		return db.Material{}, false
	}
	return m, true
}

func (s *Server) getManageableMaterial(w http.ResponseWriter, r *http.Request) (db.Material, bool) {
	m, ok := s.getAccessibleMaterial(w, r)
	if !ok {
		return db.Material{}, false
	}
	if !s.canManage(r, currentUser(r.Context()), m) {
		writeError(w, http.StatusForbidden, "нет прав на изменение")
		return db.Material{}, false
	}
	return m, true
}

func (s *Server) handleGetMaterial(w http.ResponseWriter, r *http.Request) {
	m, ok := s.getAccessibleMaterial(w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, m)
}

func (s *Server) handleUpdateMaterial(w http.ResponseWriter, r *http.Request) {
	m, ok := s.getManageableMaterial(w, r)
	if !ok {
		return
	}
	var req struct {
		Title         string  `json:"title"`
		Description   string  `json:"description"`
		Condition     string  `json:"condition"`
		Location      string  `json:"location"`
		OriginDate    string  `json:"origin_date"`
		Content       *string `json:"content"`
		ContentFormat *string `json:"content_format"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "неверный формат запроса")
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		writeError(w, http.StatusBadRequest, "укажите название материала")
		return
	}
	content := m.Content
	contentFormat := m.ContentFormat
	if req.Content != nil {
		content = req.Content
	}
	if req.ContentFormat != nil {
		contentFormat = req.ContentFormat
	}
	err := s.q.UpdateMaterial(r.Context(), db.UpdateMaterialParams{
		ID:            m.ID,
		Title:         strings.TrimSpace(req.Title),
		Description:   req.Description,
		Condition:     req.Condition,
		Location:      req.Location,
		OriginDate:    parseDate(req.OriginDate),
		Content:       content,
		ContentFormat: contentFormat,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось обновить материал")
		return
	}
	updated, err := s.q.GetMaterial(r.Context(), m.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось прочитать материал")
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (s *Server) handleDeleteMaterial(w http.ResponseWriter, r *http.Request) {
	m, ok := s.getAccessibleMaterial(w, r)
	if !ok {
		return
	}
	user := currentUser(r.Context())
	if !s.canManage(r, user, m) {
		writeError(w, http.StatusForbidden, "нет прав на удаление")
		return
	}
	if err := s.q.DeleteMaterial(r.Context(), m.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось удалить материал")
		return
	}
	for _, f := range m.Files {
		s.dropObject(r.Context(), f.Key)
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) setStatus(w http.ResponseWriter, r *http.Request, status string) {
	m, ok := s.getAccessibleMaterial(w, r)
	if !ok {
		return
	}
	user := currentUser(r.Context())
	if !isModerator(user) {
		writeError(w, http.StatusForbidden, "модерация доступна библиотекарю и завучу")
		return
	}
	if user.Role == db.RoleLibrarian && m.Category == db.CategoryDocuments {
		hasPerm, err := s.q.HasPermission(r.Context(), m.ID, user.ID)
		if err != nil || !hasPerm {
			writeError(w, http.StatusForbidden, "документы модерирует завуч")
			return
		}
	}
	if err := s.q.SetMaterialStatus(r.Context(), m.ID, status); err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось изменить статус")
		return
	}
	updated, err := s.q.GetMaterial(r.Context(), m.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось прочитать материал")
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (s *Server) handleApproveMaterial(w http.ResponseWriter, r *http.Request) {
	s.setStatus(w, r, db.StatusApproved)
}

func (s *Server) handleRejectMaterial(w http.ResponseWriter, r *http.Request) {
	s.setStatus(w, r, db.StatusRejected)
}

func (s *Server) handleMaterialQR(w http.ResponseWriter, r *http.Request) {
	m, ok := s.getAccessibleMaterial(w, r)
	if !ok {
		return
	}
	link := strings.TrimRight(s.cfg.PublicURL, "/") + fmt.Sprintf("/material/%d", m.ID)
	png, err := qrcode.Encode(link, qrcode.Medium, 512)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось создать QR-код")
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=material-%d-qr.png", m.ID))
	_, _ = w.Write(png)
}
