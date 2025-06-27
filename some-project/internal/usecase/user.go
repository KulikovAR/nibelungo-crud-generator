package usecase

import (
	"context"
	"errors"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/internal/domain"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/internal/repository"
)

type userUseCase interface {
	Create(ctx context.Context, entity *domain.user) error
	Get(ctx context.Context, id string) (*domain.user, error)
	Update(ctx context.Context, entity *domain.user) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.user, error)
}

type userUseCase struct {
	repo repository.userRepository
}

func NewuserUseCase(repo repository.userRepository) userUseCase {
	return &userUseCase{repo: repo}
}

func (uc *userUseCase) Create(ctx context.Context, entity *domain.user) error {
	if entity.ID == "" {
		entity = domain.Newuser()
	}
	return uc.repo.Create(ctx, entity)
}

func (uc *userUseCase) Get(ctx context.Context, id string) (*domain.user, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return uc.repo.Get(ctx, id)
}

func (uc *userUseCase) Update(ctx context.Context, entity *domain.user) error {
	if entity.ID == "" {
		return errors.New("id is required")
	}
	return uc.repo.Update(ctx, entity)
}

func (uc *userUseCase) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *userUseCase) List(ctx context.Context) ([]*domain.user, error) {
	return uc.repo.List(ctx)
}
