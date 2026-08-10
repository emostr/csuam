package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"csuam/backend/internal/auth"
	"csuam/backend/internal/config"
	"csuam/backend/internal/db"
	"csuam/backend/internal/httpapi"
	"csuam/backend/internal/migrate"
	"csuam/backend/internal/storage"
)

func connectDB(ctx context.Context, url string) *pgxpool.Pool {
	for i := 0; i < 30; i++ {
		pool, err := pgxpool.New(ctx, url)
		if err == nil {
			if err = pool.Ping(ctx); err == nil {
				return pool
			}
			pool.Close()
		}
		log.Printf("wait for db: %v", err)
		time.Sleep(2 * time.Second)
	}
	log.Fatal("err connect to db")
	return nil
}

func seedUsers(ctx context.Context, cfg config.Config, q *db.Queries) {
	n, err := q.CountUsers(ctx)
	if err != nil {
		log.Fatalf("err check user: %v", err)
	}
	if n > 0 {
		return
	}
	seed := []struct {
		username, password, fullName, role string
	}{
		{"zavuch", cfg.SeedHeadPass, "Администратор", db.RoleHeadTeacher},
		{"librarian", cfg.SeedLibPass, "Библиотекарь", db.RoleLibrarian},
		{"teacher", cfg.SeedTeacherPass, "Учитель", db.RoleTeacher},
	}
	for _, u := range seed {
		hash, err := auth.HashPassword(u.password)
		if err != nil {
			log.Fatalf("err hash psw: %v", err)
		}
		if _, err := q.CreateUser(ctx, u.username, hash, u.fullName, u.role); err != nil {
			log.Fatalf("err create user %s: %v", u.username, err)
		}
		log.Printf("created user %s (%s)", u.username, u.role)
	}
}

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pool := connectDB(ctx, cfg.DatabaseURL)
	defer pool.Close()

	if err := migrate.Run(ctx, pool); err != nil {
		log.Fatalf("migrations isnt appended: %v", err)
	}

	queries := db.New(pool)
	seedUsers(ctx, cfg, queries)

	store, err := storage.New(cfg)
	if err != nil {
		log.Fatalf("err minio connect: %v", err)
	}
	for i := 0; i < 30; i++ {
		if err = store.EnsureBucket(ctx); err == nil {
			break
		}
		log.Printf("minio: %v", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("minio not available: %v", err)
	}

	server := httpapi.NewServer(cfg, queries, pool, store)
	addr := ":" + cfg.Port
	log.Printf("started on %s", addr)
	if err := http.ListenAndServe(addr, server.Router()); err != nil {
		log.Fatal(err)
	}
}
