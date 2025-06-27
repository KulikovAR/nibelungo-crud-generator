package domain

import (
	"time"
	"github.com/google/uuid"
)

type user struct {
	ID int64
	Username string
	Email string
	CreatedAt string
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func Newuser() *user {
	return &user{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
