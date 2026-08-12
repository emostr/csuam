package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"csuam/backend/internal/auth"
	"csuam/backend/internal/config"
	"csuam/backend/internal/db"
	"csuam/backend/internal/httpapi"
	"csuam/backend/internal/migrate"
	"csuam/backend/internal/storage"
)

func connectDB(ctx context.Context, dsn string) *sql.DB {
	for i := 0; i < 30; i++ {
		sdb, err := sql.Open("mysql", dsn)
		if err == nil {
			if err = sdb.PingContext(ctx); err == nil {
				sdb.SetConnMaxLifetime(3 * time.Minute)
				sdb.SetMaxOpenConns(10)
				sdb.SetMaxIdleConns(10)
				return sdb
			}
			sdb.Close()
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

	sdb := connectDB(ctx, cfg.DatabaseURL)
	defer sdb.Close()

	if err := migrate.Run(ctx, sdb); err != nil {
		log.Fatalf("migrations isnt appended: %v", err)
	}

	queries := db.New(sdb)
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

	server := httpapi.NewServer(cfg, queries, store)
	addr := ":" + cfg.Port
	log.Printf("started on %s", addr)
	if err := http.ListenAndServe(addr, server.Router()); err != nil {
		log.Fatal(err)
	}
}
