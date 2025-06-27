package repository

import (
	"context"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/internal/domain"
)

type productRepository interface {
	Create(ctx context.Context, entity *domain.product) error
	Get(ctx context.Context, id string) (*domain.product, error)
	Update(ctx context.Context, entity *domain.product) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.product, error)
}
