package httpapi

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"csuam/backend/internal/db"
	"csuam/backend/internal/storage"
)

var errUnsupportedFormat = errors.New("формат файла не поддерживается")

type uploadedFile struct {
	key      string
	name     string
	mimeType string
	size     int64
}

func formFileHeaders(r *http.Request) []*multipart.FileHeader {
	if r.MultipartForm == nil {
		return nil
	}
	headers := append([]*multipart.FileHeader{}, r.MultipartForm.File["files"]...)
	return append(headers, r.MultipartForm.File["file"]...)
}

func (s *Server) uploadFile(ctx context.Context, header *multipart.FileHeader) (uploadedFile, error) {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExtensions[ext] {
		return uploadedFile{}, errUnsupportedFormat
	}
	file, err := header.Open()
	if err != nil {
		return uploadedFile{}, err
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = mime.TypeByExtension(ext)
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	key := storage.NewKey(ext)
	if err := s.store.Put(ctx, key, file, header.Size, contentType); err != nil {
		return uploadedFile{}, err
	}
	if isRasterImage(contentType) {
		if _, err := file.Seek(0, io.SeekStart); err == nil {
			if thumb, err := makeThumbnail(file); err == nil {
				_ = s.store.Put(ctx, thumbKey(key), bytes.NewReader(thumb), int64(len(thumb)), "image/jpeg")
			}
		}
	}
	return uploadedFile{key: key, name: header.Filename, mimeType: contentType, size: header.Size}, nil
}

func (s *Server) uploadFiles(ctx context.Context, headers []*multipart.FileHeader) ([]uploadedFile, error) {
	uploads := make([]uploadedFile, 0, len(headers))
	for _, header := range headers {
		up, err := s.uploadFile(ctx, header)
		if err != nil {
			s.dropUploads(ctx, uploads)
			return nil, err
		}
		uploads = append(uploads, up)
	}
	return uploads, nil
}

func (s *Server) dropUploads(ctx context.Context, uploads []uploadedFile) {
	for _, up := range uploads {
		s.dropObject(ctx, up.key)
	}
}

func (s *Server) dropObject(ctx context.Context, key string) {
	_ = s.store.Delete(ctx, key)
	_ = s.store.Delete(ctx, thumbKey(key))
}

func (s *Server) attachFiles(ctx context.Context, materialID int64, uploads []uploadedFile) error {
	for _, up := range uploads {
		_, err := s.q.AddMaterialFile(ctx, db.AddMaterialFileParams{
			MaterialID: materialID,
			Key:        up.key,
			Name:       up.name,
			Mime:       up.mimeType,
			Size:       up.size,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) handleAddMaterialFiles(w http.ResponseWriter, r *http.Request) {
	m, ok := s.getManageableMaterial(w, r)
	if !ok {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "не удалось прочитать форму (слишком большой файл?)")
		return
	}
	headers := formFileHeaders(r)
	if len(headers) == 0 {
		writeError(w, http.StatusBadRequest, "выберите хотя бы один файл")
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
	if err := s.attachFiles(r.Context(), m.ID, uploads); err != nil {
		s.dropUploads(r.Context(), uploads)
		writeError(w, http.StatusInternalServerError, "не удалось прикрепить файлы")
		return
	}
	updated, err := s.q.GetMaterial(r.Context(), m.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "файлы добавлены, но не удалось прочитать материал")
		return
	}
	writeJSON(w, http.StatusCreated, updated)
}

func (s *Server) handleDeleteMaterialFile(w http.ResponseWriter, r *http.Request) {
	m, ok := s.getManageableMaterial(w, r)
	if !ok {
		return
	}
	fileID, err := urlID(r, "fileId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "неверный идентификатор файла")
		return
	}
	file, ok := findFile(m, fileID)
	if !ok {
		writeError(w, http.StatusNotFound, "файл не найден")
		return
	}
	if err := s.q.DeleteMaterialFile(r.Context(), m.ID, file.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось удалить файл")
		return
	}
	s.dropObject(r.Context(), file.Key)
	updated, err := s.q.GetMaterial(r.Context(), m.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "файл удалён, но не удалось прочитать материал")
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func findFile(m db.Material, fileID int64) (db.MaterialFile, bool) {
	for _, f := range m.Files {
		if f.ID == fileID {
			return f, true
		}
	}
	return db.MaterialFile{}, false
}

func (s *Server) handleMaterialFile(w http.ResponseWriter, r *http.Request) {
	m, ok := s.getAccessibleMaterial(w, r)
	if !ok {
		return
	}
	if len(m.Files) == 0 {
		writeError(w, http.StatusNotFound, "у материала нет файлов")
		return
	}
	s.serveFile(w, r, m.Files[0])
}

func (s *Server) handleMaterialFileByID(w http.ResponseWriter, r *http.Request) {
	m, ok := s.getAccessibleMaterial(w, r)
	if !ok {
		return
	}
	fileID, err := urlID(r, "fileId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "неверный идентификатор файла")
		return
	}
	file, ok := findFile(m, fileID)
	if !ok {
		writeError(w, http.StatusNotFound, "файл не найден")
		return
	}
	s.serveFile(w, r, file)
}

func (s *Server) serveFile(w http.ResponseWriter, r *http.Request, f db.MaterialFile) {
	if r.URL.Query().Get("thumb") == "1" && isRasterImage(f.Mime) {
		if obj, err := s.store.Get(r.Context(), thumbKey(f.Key)); err == nil {
			if stat, err := obj.Stat(); err == nil {
				defer obj.Close()
				w.Header().Set("Content-Type", "image/jpeg")
				w.Header().Set("Cache-Control", "private, max-age=86400")
				http.ServeContent(w, r, "thumb.jpg", stat.LastModified, obj)
				return
			}
			obj.Close()
		}
	}
	obj, err := s.store.Get(r.Context(), f.Key)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось открыть файл")
		return
	}
	defer obj.Close()
	stat, err := obj.Stat()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "файл недоступен в хранилище")
		return
	}
	disposition := "inline"
	if r.URL.Query().Get("download") == "1" {
		disposition = "attachment"
	}
	w.Header().Set("Content-Type", f.Mime)
	w.Header().Set("Cache-Control", "private, max-age=86400")
	w.Header().Set("Content-Disposition", fmt.Sprintf("%s; filename*=UTF-8''%s", disposition, url.PathEscape(f.Name)))
	http.ServeContent(w, r, f.Name, stat.LastModified, obj)
}
