package db

import (
	"time"
)

const (
	RoleTeacher     = "teacher"
	RoleLibrarian   = "librarian"
	RoleHeadTeacher = "head_teacher"
)

const (
	CategoryAwards    = "awards"
	CategoryPhotos    = "photos"
	CategoryVideos    = "videos"
	CategoryLibrary   = "library"
	CategoryDocuments = "documents"
)

const (
	StatusPending  = "pending"
	StatusApproved = "approved"
	StatusRejected = "rejected"
)

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type Material struct {
	ID            int64          `json:"id"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	Category      string         `json:"category"`
	Status        string         `json:"status"`
	Content       *string        `json:"content"`
	ContentFormat *string        `json:"content_format"`
	Condition     string         `json:"condition"`
	Location      string         `json:"location"`
	OriginDate    *time.Time     `json:"origin_date"`
	CreatedBy     *int64         `json:"created_by"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	AuthorName    *string        `json:"author_name"`
	Files         []MaterialFile `json:"files"`
}

type MaterialFile struct {
	ID         int64     `json:"id"`
	MaterialID int64     `json:"material_id"`
	Key        string    `json:"-"`
	Name       string    `json:"name"`
	Mime       string    `json:"mime"`
	Size       int64     `json:"size"`
	CreatedAt  time.Time `json:"created_at"`
}

type Permission struct {
	ID         int64     `json:"id"`
	MaterialID int64     `json:"material_id"`
	UserID     int64     `json:"user_id"`
	FullName   string    `json:"full_name"`
	Username   string    `json:"username"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
}

type Loan struct {
	ID            int64      `json:"id"`
	MaterialID    int64      `json:"material_id"`
	MaterialTitle string     `json:"material_title"`
	BorrowerName  string     `json:"borrower_name"`
	Note          string     `json:"note"`
	TakenAt       time.Time  `json:"taken_at"`
	DueDate       time.Time  `json:"due_date"`
	ReturnedAt    *time.Time `json:"returned_at"`
}

type CategoryCount struct {
	Category string `json:"category"`
	Count    int64  `json:"count"`
}

type MonthCount struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}
