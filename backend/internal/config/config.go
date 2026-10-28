package config

import (
	"os"
)

type Config struct {
	Port            string
	DatabaseURL     string
	JWTSecret       string
	PublicURL       string
	S3Endpoint      string
	S3AccessKey     string
	S3SecretKey     string
	S3Bucket        string
	S3Region        string
	S3UseSSL        bool
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
		DatabaseURL:     env("DATABASE_URL", "postgres://csuam:csuam@localhost:5432/csuam?sslmode=disable"),
		JWTSecret:       env("JWT_SECRET", "dev-secret-change-me"),
		PublicURL:       env("PUBLIC_URL", "http://localhost:5173"),
		S3Endpoint:      env("S3_ENDPOINT", "localhost:3900"),
		S3AccessKey:     env("S3_ACCESS_KEY", "GKcsuamdevaccesskey"),
		S3SecretKey:     env("S3_SECRET_KEY", "csuam-dev-secret-change-me"),
		S3Bucket:        env("S3_BUCKET", "csuam"),
		S3Region:        env("S3_REGION", "garage"),
		S3UseSSL:        env("S3_USE_SSL", "false") == "true",
		SeedHeadPass:    env("SEED_HEAD_TEACHER_PASSWORD", "zavuch123"),
		SeedLibPass:     env("SEED_LIBRARIAN_PASSWORD", "biblio123"),
		SeedTeacherPass: env("SEED_TEACHER_PASSWORD", "teacher123"),
	}
}
