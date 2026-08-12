package config

import (
	"os"
)

type Config struct {
	Port            string
	DatabaseURL     string
	JWTSecret       string
	PublicURL       string
	MinioEndpoint   string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioBucket     string
	MinioUseSSL     bool
	SeedHeadPass    string
	SeedLibPass     string
	SeedTeacherPass string
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func Load() Config {
	return Config{
		Port:            env("PORT", "8080"),
		DatabaseURL:     env("DATABASE_URL", "csuam:csuam@tcp(localhost:3306)/csuam?parseTime=true&multiStatements=true"),
		JWTSecret:       env("JWT_SECRET", "dev-secret-change-me"),
		PublicURL:       env("PUBLIC_URL", "http://localhost:5173"),
		MinioEndpoint:   env("MINIO_ENDPOINT", "localhost:9000"),
		MinioAccessKey:  env("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecretKey:  env("MINIO_SECRET_KEY", "minioadmin"),
		MinioBucket:     env("MINIO_BUCKET", "csuam"),
		MinioUseSSL:     env("MINIO_USE_SSL", "false") == "true",
		SeedHeadPass:    env("SEED_HEAD_TEACHER_PASSWORD", "zavuch123"),
		SeedLibPass:     env("SEED_LIBRARIAN_PASSWORD", "biblio123"),
		SeedTeacherPass: env("SEED_TEACHER_PASSWORD", "teacher123"),
	}
}
