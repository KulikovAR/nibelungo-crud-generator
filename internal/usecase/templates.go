package usecase

const (
	// Шаблоны для domain
	domainEntityTemplate = `package domain

type {{.Name}} struct {
	{{- range .Fields}}
	{{.Name}} {{.Type}} {{- range .Tags}} {{.}}{{- end}}
	{{- end}}
}

func New{{.Name}}() *{{.Name}} {
	return &{{.Name}}{}
}
`

	// Шаблоны для repository
	repositoryInterfaceTemplate = `package repository

import (
	"context"
	"github.com/KulikovAR/nibelungo-crud-generator/internal/domain"
)

type {{.Name}}Repository interface {
	Create(ctx context.Context, entity *domain.{{.Name}}) error
	Get(ctx context.Context, id string) (*domain.{{.Name}}, error)
	Update(ctx context.Context, entity *domain.{{.Name}}) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.{{.Name}}, error)
}
`

	// Шаблоны для postgres repository
	postgresRepositoryTemplate = `package postgres

import (
	"context"
	"database/sql"
	"github.com/KulikovAR/nibelungo-crud-generator/internal/domain"
	"github.com/KulikovAR/nibelungo-crud-generator/internal/repository"
)

type {{.Name}}Repository struct {
	db *sql.DB
}

func New{{.Name}}Repository(db *sql.DB) repository.{{.Name}}Repository {
	return &{{.Name}}Repository{db: db}
}

func (r *{{.Name}}Repository) Create(ctx context.Context, entity *domain.{{.Name}}) error {
	// TODO: Implement
	return nil
}

func (r *{{.Name}}Repository) Get(ctx context.Context, id string) (*domain.{{.Name}}, error) {
	// TODO: Implement
	return nil, nil
}

func (r *{{.Name}}Repository) Update(ctx context.Context, entity *domain.{{.Name}}) error {
	// TODO: Implement
	return nil
}

func (r *{{.Name}}Repository) Delete(ctx context.Context, id string) error {
	// TODO: Implement
	return nil
}

func (r *{{.Name}}Repository) List(ctx context.Context) ([]*domain.{{.Name}}, error) {
	// TODO: Implement
	return nil, nil
}
`

	// Шаблоны для mongodb repository
	mongodbRepositoryTemplate = `package mongodb

import (
	"context"
	"github.com/KulikovAR/nibelungo-crud-generator/internal/domain"
	"github.com/KulikovAR/nibelungo-crud-generator/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type {{.Name}}Repository struct {
	collection *mongo.Collection
}

func New{{.Name}}Repository(collection *mongo.Collection) repository.{{.Name}}Repository {
	return &{{.Name}}Repository{collection: collection}
}

func (r *{{.Name}}Repository) Create(ctx context.Context, entity *domain.{{.Name}}) error {
	// TODO: Implement
	return nil
}

func (r *{{.Name}}Repository) Get(ctx context.Context, id string) (*domain.{{.Name}}, error) {
	// TODO: Implement
	return nil, nil
}

func (r *{{.Name}}Repository) Update(ctx context.Context, entity *domain.{{.Name}}) error {
	// TODO: Implement
	return nil
}

func (r *{{.Name}}Repository) Delete(ctx context.Context, id string) error {
	// TODO: Implement
	return nil
}

func (r *{{.Name}}Repository) List(ctx context.Context) ([]*domain.{{.Name}}, error) {
	// TODO: Implement
	return nil, nil
}
`

	// Шаблоны для usecase
	usecaseTemplate = `package usecase

import (
	"context"
	"github.com/KulikovAR/nibelungo-crud-generator/internal/domain"
	"github.com/KulikovAR/nibelungo-crud-generator/internal/repository"
)

type {{.Name}}UseCase interface {
	Create(ctx context.Context, entity *domain.{{.Name}}) error
	Get(ctx context.Context, id string) (*domain.{{.Name}}, error)
	Update(ctx context.Context, entity *domain.{{.Name}}) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.{{.Name}}, error)
}

type {{.Name | ToLower}}UseCase struct {
	repo repository.{{.Name}}Repository
}

func New{{.Name}}UseCase(repo repository.{{.Name}}Repository) {{.Name}}UseCase {
	return &{{.Name | ToLower}}UseCase{repo: repo}
}

func (uc *{{.Name | ToLower}}UseCase) Create(ctx context.Context, entity *domain.{{.Name}}) error {
	return uc.repo.Create(ctx, entity)
}

func (uc *{{.Name | ToLower}}UseCase) Get(ctx context.Context, id string) (*domain.{{.Name}}, error) {
	return uc.repo.Get(ctx, id)
}

func (uc *{{.Name | ToLower}}UseCase) Update(ctx context.Context, entity *domain.{{.Name}}) error {
	return uc.repo.Update(ctx, entity)
}

func (uc *{{.Name | ToLower}}UseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *{{.Name | ToLower}}UseCase) List(ctx context.Context) ([]*domain.{{.Name}}, error) {
	return uc.repo.List(ctx)
}
`

	// Шаблоны для controller
	controllerTemplate = `package controller

import (
	"context"
	"github.com/KulikovAR/nibelungo-crud-generator/internal/domain"
	"github.com/KulikovAR/nibelungo-crud-generator/internal/usecase"
)

type {{.Name}}Controller struct {
	useCase usecase.{{.Name}}UseCase
}

func New{{.Name}}Controller(useCase usecase.{{.Name}}UseCase) *{{.Name}}Controller {
	return &{{.Name}}Controller{useCase: useCase}
}

func (c *{{.Name}}Controller) Create(ctx context.Context, entity *domain.{{.Name}}) error {
	return c.useCase.Create(ctx, entity)
}

func (c *{{.Name}}Controller) Get(ctx context.Context, id string) (*domain.{{.Name}}, error) {
	return c.useCase.Get(ctx, id)
}

func (c *{{.Name}}Controller) Update(ctx context.Context, entity *domain.{{.Name}}) error {
	return c.useCase.Update(ctx, entity)
}

func (c *{{.Name}}Controller) Delete(ctx context.Context, id string) error {
	return c.useCase.Delete(ctx, id)
}

func (c *{{.Name}}Controller) List(ctx context.Context) ([]*domain.{{.Name}}, error) {
	return c.useCase.List(ctx)
}
`

	// Шаблоны для миграций
	postgresMigrationTemplate = `-- +migrate Up
CREATE TABLE {{.Name | ToSnakeCase}}s (
    id VARCHAR(36) PRIMARY KEY,
    {{- range .Fields}}
    {{.Name | ToSnakeCase}} {{.Type | ToPostgresType}},
    {{- end}}
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE {{.Name | ToSnakeCase}}s;
`

	mongodbMigrationTemplate = `{
    "collection": "{{.Name | ToSnakeCase}}s",
    "indexes": [
        {
            "keys": {
                "id": 1
            },
            "options": {
                "unique": true
            }
        }
    ]
}`
)
