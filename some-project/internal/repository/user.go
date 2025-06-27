package repository

import (
	"context"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/internal/domain"
)

type userRepository interface {
	Create(ctx context.Context, entity *domain.user) error
	Get(ctx context.Context, id string) (*domain.user, error)
	Update(ctx context.Context, entity *domain.user) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.user, error)
}
