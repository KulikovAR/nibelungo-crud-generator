package usecase

import (
	"context"
	"errors"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/internal/domain"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/internal/repository"
)

type productUseCase interface {
	Create(ctx context.Context, entity *domain.product) error
	Get(ctx context.Context, id string) (*domain.product, error)
	Update(ctx context.Context, entity *domain.product) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.product, error)
}

type productUseCase struct {
	repo repository.productRepository
}

func NewproductUseCase(repo repository.productRepository) productUseCase {
	return &productUseCase{repo: repo}
}

func (uc *productUseCase) Create(ctx context.Context, entity *domain.product) error {
	if entity.ID == "" {
		entity = domain.Newproduct()
	}
	return uc.repo.Create(ctx, entity)
}

func (uc *productUseCase) Get(ctx context.Context, id string) (*domain.product, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return uc.repo.Get(ctx, id)
}

func (uc *productUseCase) Update(ctx context.Context, entity *domain.product) error {
	if entity.ID == "" {
		return errors.New("id is required")
	}
	return uc.repo.Update(ctx, entity)
}

func (uc *productUseCase) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *productUseCase) List(ctx context.Context) ([]*domain.product, error) {
	return uc.repo.List(ctx)
}
