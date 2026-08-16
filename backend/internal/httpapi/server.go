package httpapi

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"csuam/backend/internal/config"
	"csuam/backend/internal/db"
	"csuam/backend/internal/storage"
)

type Server struct {
	cfg   config.Config
	q     *db.Queries
	store *storage.Storage
}

func NewServer(cfg config.Config, q *db.Queries, store *storage.Storage) *Server {
	return &Server{cfg: cfg, q: q, store: store}
}

func (s *Server) Router() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Post("/auth/login", s.handleLogin)
		r.Post("/auth/logout", s.handleLogout)

		r.Group(func(r chi.Router) {
			r.Use(s.authMiddleware)

			r.Get("/auth/me", s.handleMe)
			r.Post("/auth/password", s.handleChangePassword)
			r.Get("/notifications", s.handleNotifications)

			r.Route("/users", func(r chi.Router) {
				r.Get("/", s.handleListUsers)
				r.Post("/", s.handleCreateUser)
				r.Delete("/{id}", s.handleDeleteUser)
				r.Post("/{id}/password", s.handleResetPassword)
			})

			r.Route("/materials", func(r chi.Router) {
				r.Get("/", s.handleListMaterials)
				r.Post("/", s.handleCreateMaterial)
				r.Post("/import", s.handleImportMaterial)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", s.handleGetMaterial)
					r.Put("/", s.handleUpdateMaterial)
					r.Delete("/", s.handleDeleteMaterial)
					r.Post("/approve", s.handleApproveMaterial)
					r.Post("/reject", s.handleRejectMaterial)
					r.Get("/file", s.handleMaterialFile)
					r.Post("/files", s.handleAddMaterialFiles)
					r.Get("/files/{fileId}", s.handleMaterialFileByID)
					r.Delete("/files/{fileId}", s.handleDeleteMaterialFile)
					r.Get("/qr", s.handleMaterialQR)
					r.Get("/export", s.handleExportMaterial)
					r.Get("/permissions", s.handleListPermissions)
					r.Post("/permissions", s.handleGrantPermission)
					r.Delete("/permissions/{userId}", s.handleRevokePermission)
				})
			})

			r.Route("/loans", func(r chi.Router) {
				r.Get("/", s.handleListLoans)
				r.Post("/", s.handleCreateLoan)
				r.Post("/{id}/return", s.handleReturnLoan)
			})

			r.Get("/analytics/summary", s.handleAnalyticsSummary)
		})
	})

	return r
}
