package domain

import (
	"time"
	"github.com/google/uuid"
)

type product struct {
	ID int64
	Name string
	Description string
	Price float64
	CreatedAt string
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func Newproduct() *product {
	return &product{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
