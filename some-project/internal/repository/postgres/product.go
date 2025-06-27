package postgres

import (
	"context"
	"database/sql"
	"time"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/internal/domain"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/internal/repository"
)

type productRepository struct {
	db *sql.DB
}

func NewproductRepository(db *sql.DB) repository.productRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, entity *domain.product) error {
	query := `
		INSERT INTO products (
			id, id, name, description, price, created_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`
	
	_, err := r.db.ExecContext(ctx, query, 
		entity.ID, 
		entity.ID, entity.Name, entity.Description, entity.Price, entity.CreatedAt, 
		entity.CreatedAt, 
		entity.UpdatedAt,
	)
	return err
}

func (r *productRepository) Get(ctx context.Context, id string) (*domain.product, error) {
	query := `SELECT id, id, name, description, price, created_at, created_at, updated_at FROM products WHERE id = $1`
	
	var entity domain.product
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&entity.ID,
		&entity.ID, &entity.Name, &entity.Description, &entity.Price, &entity.CreatedAt, 
		&entity.CreatedAt,
		&entity.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *productRepository) Update(ctx context.Context, entity *domain.product) error {
	entity.UpdatedAt = time.Now()
	query := `
		UPDATE products SET 
			id = $2, name = $3, description = $4, price = $5, created_at = $6, 
			updated_at = $7
		WHERE id = $1
	`
	
	_, err := r.db.ExecContext(ctx, query, 
		entity.ID,
		entity.ID, entity.Name, entity.Description, entity.Price, entity.CreatedAt, 
		entity.UpdatedAt,
	)
	return err
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *productRepository) List(ctx context.Context) ([]*domain.product, error) {
	query := `SELECT id, id, name, description, price, created_at, created_at, updated_at FROM products ORDER BY created_at DESC`
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var entities []*domain.product
	for rows.Next() {
		var entity domain.product
		err := rows.Scan(
			&entity.ID,
			&entity.ID, &entity.Name, &entity.Description, &entity.Price, &entity.CreatedAt, 
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
