package postgres

import (
	"context"
	"database/sql"
	"time"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/internal/domain"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/internal/repository"
)

type userRepository struct {
	db *sql.DB
}

func NewuserRepository(db *sql.DB) repository.userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, entity *domain.user) error {
	query := `
		INSERT INTO users (
			id, id, username, email, created_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
	`
	
	_, err := r.db.ExecContext(ctx, query, 
		entity.ID, 
		entity.ID, entity.Username, entity.Email, entity.CreatedAt, 
		entity.CreatedAt, 
		entity.UpdatedAt,
	)
	return err
}

func (r *userRepository) Get(ctx context.Context, id string) (*domain.user, error) {
	query := `SELECT id, id, username, email, created_at, created_at, updated_at FROM users WHERE id = $1`
	
	var entity domain.user
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&entity.ID,
		&entity.ID, &entity.Username, &entity.Email, &entity.CreatedAt, 
		&entity.CreatedAt,
		&entity.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *userRepository) Update(ctx context.Context, entity *domain.user) error {
	entity.UpdatedAt = time.Now()
	query := `
		UPDATE users SET 
			id = $2, username = $3, email = $4, created_at = $5, 
			updated_at = $6
		WHERE id = $1
	`
	
	_, err := r.db.ExecContext(ctx, query, 
		entity.ID,
		entity.ID, entity.Username, entity.Email, entity.CreatedAt, 
		entity.UpdatedAt,
	)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *userRepository) List(ctx context.Context) ([]*domain.user, error) {
	query := `SELECT id, id, username, email, created_at, created_at, updated_at FROM users ORDER BY created_at DESC`
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var entities []*domain.user
	for rows.Next() {
		var entity domain.user
		err := rows.Scan(
			&entity.ID,
			&entity.ID, &entity.Username, &entity.Email, &entity.CreatedAt, 
			&entity.CreatedAt,
			&entity.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entities = append(entities, &entity)
	}
	return entities, nil
}
