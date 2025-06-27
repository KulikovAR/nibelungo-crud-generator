package usecase

const (
	// Шаблоны для domain
	domainEntityTemplate = `package domain

import (
	"time"
	"github.com/google/uuid"
)

type {{.Name}} struct {
	{{- range .Fields}}
	{{.Name}} {{.Type}} {{- if .Tags}} ` + "`" + `{{range $i, $tag := .Tags}}{{if $i}} {{end}}{{$tag}}{{end}}` + "`" + `{{- end}}
	{{- end}}
	CreatedAt time.Time ` + "`" + `json:"created_at" db:"created_at"` + "`" + `
	UpdatedAt time.Time ` + "`" + `json:"updated_at" db:"updated_at"` + "`" + `
}

func New{{.Name}}() *{{.Name}} {
	return &{{.Name}}{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
`

	// Шаблоны для repository
	repositoryInterfaceTemplate = `package repository

import (
	"context"
	"github.com/KulikovAR/{{.Module}}/internal/domain"
)

type {{.Entity.Name}}Repository interface {
	Create(ctx context.Context, entity *domain.{{.Entity.Name}}) error
	Get(ctx context.Context, id string) (*domain.{{.Entity.Name}}, error)
	Update(ctx context.Context, entity *domain.{{.Entity.Name}}) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.{{.Entity.Name}}, error)
}
`

	// Шаблоны для postgres repository
	postgresRepositoryTemplate = `package postgres

import (
	"context"
	"database/sql"
	"time"
	"github.com/KulikovAR/{{.Module}}/internal/domain"
	"github.com/KulikovAR/{{.Module}}/internal/repository"
)

type {{.Entity.Name}}Repository struct {
	db *sql.DB
}

func New{{.Entity.Name}}Repository(db *sql.DB) repository.{{.Entity.Name}}Repository {
	return &{{.Entity.Name}}Repository{db: db}
}

func (r *{{.Entity.Name}}Repository) Create(ctx context.Context, entity *domain.{{.Entity.Name}}) error {
	query := ` + "`" + `
		INSERT INTO {{.Entity.Name | ToSnakeCase}}s (
			id, {{range $i, $field := .Entity.Fields}}{{if $i}}, {{end}}{{$field.Name | ToSnakeCase}}{{end}}, created_at, updated_at
		) VALUES (
			$1, {{range $i, $field := .Entity.Fields}}{{if $i}}, {{end}}${{add $i 2}}{{end}}, ${{add (len .Entity.Fields) 2}}, ${{add (len .Entity.Fields) 3}}
		)
	` + "`" + `
	
	_, err := r.db.ExecContext(ctx, query, 
		entity.ID, 
		{{range $i, $field := .Entity.Fields}}entity.{{$field.Name}}, {{end}}
		entity.CreatedAt, 
		entity.UpdatedAt,
	)
	return err
}

func (r *{{.Entity.Name}}Repository) Get(ctx context.Context, id string) (*domain.{{.Entity.Name}}, error) {
	query := ` + "`" + `SELECT id, {{range $i, $field := .Entity.Fields}}{{if $i}}, {{end}}{{$field.Name | ToSnakeCase}}{{end}}, created_at, updated_at FROM {{.Entity.Name | ToSnakeCase}}s WHERE id = $1` + "`" + `
	
	var entity domain.{{.Entity.Name}}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&entity.ID,
		{{range $i, $field := .Entity.Fields}}&entity.{{$field.Name}}, {{end}}
		&entity.CreatedAt,
		&entity.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *{{.Entity.Name}}Repository) Update(ctx context.Context, entity *domain.{{.Entity.Name}}) error {
	entity.UpdatedAt = time.Now()
	query := ` + "`" + `
		UPDATE {{.Entity.Name | ToSnakeCase}}s SET 
			{{range $i, $field := .Entity.Fields}}{{if $i}}, {{end}}{{$field.Name | ToSnakeCase}} = ${{add $i 2}}{{end}}, 
			updated_at = ${{add (len .Entity.Fields) 2}}
		WHERE id = $1
	` + "`" + `
	
	_, err := r.db.ExecContext(ctx, query, 
		entity.ID,
		{{range $i, $field := .Entity.Fields}}entity.{{$field.Name}}, {{end}}
		entity.UpdatedAt,
	)
	return err
}

func (r *{{.Entity.Name}}Repository) Delete(ctx context.Context, id string) error {
	query := ` + "`" + `DELETE FROM {{.Entity.Name | ToSnakeCase}}s WHERE id = $1` + "`" + `
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *{{.Entity.Name}}Repository) List(ctx context.Context) ([]*domain.{{.Entity.Name}}, error) {
	query := ` + "`" + `SELECT id, {{range $i, $field := .Entity.Fields}}{{if $i}}, {{end}}{{$field.Name | ToSnakeCase}}{{end}}, created_at, updated_at FROM {{.Entity.Name | ToSnakeCase}}s ORDER BY created_at DESC` + "`" + `
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var entities []*domain.{{.Entity.Name}}
	for rows.Next() {
		var entity domain.{{.Entity.Name}}
		err := rows.Scan(
			&entity.ID,
			{{range $i, $field := .Entity.Fields}}&entity.{{$field.Name}}, {{end}}
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
`

	// Шаблоны для mongodb repository
	mongodbRepositoryTemplate = `package mongodb

import (
	"context"
	"time"
	"github.com/KulikovAR/{{.Module}}/internal/domain"
	"github.com/KulikovAR/{{.Module}}/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type {{.Entity.Name}}Repository struct {
	collection *mongo.Collection
}

func New{{.Entity.Name}}Repository(collection *mongo.Collection) repository.{{.Entity.Name}}Repository {
	return &{{.Entity.Name}}Repository{collection: collection}
}

func (r *{{.Entity.Name}}Repository) Create(ctx context.Context, entity *domain.{{.Entity.Name}}) error {
	entity.CreatedAt = time.Now()
	entity.UpdatedAt = time.Now()
	
	_, err := r.collection.InsertOne(ctx, entity)
	return err
}

func (r *{{.Entity.Name}}Repository) Get(ctx context.Context, id string) (*domain.{{.Entity.Name}}, error) {
	var entity domain.{{.Entity.Name}}
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&entity)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *{{.Entity.Name}}Repository) Update(ctx context.Context, entity *domain.{{.Entity.Name}}) error {
	entity.UpdatedAt = time.Now()
	
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"id": entity.ID},
		bson.M{"$set": entity},
	)
	return err
}

func (r *{{.Entity.Name}}Repository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (r *{{.Entity.Name}}Repository) List(ctx context.Context) ([]*domain.{{.Entity.Name}}, error) {
	opts := options.Find().SetSort({{print "bson.D{{\"created_at\", -1}}"}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var entities []*domain.{{.Entity.Name}}
	if err = cursor.All(ctx, &entities); err != nil {
		return nil, err
	}
	return entities, nil
}
`

	// Шаблоны для usecase
	usecaseTemplate = `package usecase

import (
	"context"
	"errors"
	"github.com/KulikovAR/{{.Module}}/internal/domain"
	"github.com/KulikovAR/{{.Module}}/internal/repository"
)

type {{.Entity.Name}}UseCase interface {
	Create(ctx context.Context, entity *domain.{{.Entity.Name}}) error
	Get(ctx context.Context, id string) (*domain.{{.Entity.Name}}, error)
	Update(ctx context.Context, entity *domain.{{.Entity.Name}}) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.{{.Entity.Name}}, error)
}

type {{.Entity.Name | ToLower}}UseCase struct {
	repo repository.{{.Entity.Name}}Repository
}

func New{{.Entity.Name}}UseCase(repo repository.{{.Entity.Name}}Repository) {{.Entity.Name}}UseCase {
	return &{{.Entity.Name | ToLower}}UseCase{repo: repo}
}

func (uc *{{.Entity.Name | ToLower}}UseCase) Create(ctx context.Context, entity *domain.{{.Entity.Name}}) error {
	if entity.ID == "" {
		entity = domain.New{{.Entity.Name}}()
	}
	return uc.repo.Create(ctx, entity)
}

func (uc *{{.Entity.Name | ToLower}}UseCase) Get(ctx context.Context, id string) (*domain.{{.Entity.Name}}, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return uc.repo.Get(ctx, id)
}

func (uc *{{.Entity.Name | ToLower}}UseCase) Update(ctx context.Context, entity *domain.{{.Entity.Name}}) error {
	if entity.ID == "" {
		return errors.New("id is required")
	}
	return uc.repo.Update(ctx, entity)
}

func (uc *{{.Entity.Name | ToLower}}UseCase) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *{{.Entity.Name | ToLower}}UseCase) List(ctx context.Context) ([]*domain.{{.Entity.Name}}, error) {
	return uc.repo.List(ctx)
}
`

	// Шаблоны для REST controller
	restControllerTemplate = `package controller

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/KulikovAR/{{.Module}}/internal/domain"
	"github.com/KulikovAR/{{.Module}}/internal/usecase"
)

type {{.Entity.Name}}Controller struct {
	useCase usecase.{{.Entity.Name}}UseCase
}

func New{{.Entity.Name}}Controller(useCase usecase.{{.Entity.Name}}UseCase) *{{.Entity.Name}}Controller {
	return &{{.Entity.Name}}Controller{useCase: useCase}
}

// Create{{.Entity.Name}} godoc
// @Summary Create a new {{.Entity.Name | ToLower}}
// @Description Create a new {{.Entity.Name | ToLower}} with the input payload
// @Tags {{.Entity.Name | ToLower}}s
// @Accept json
// @Produce json
// @Param {{.Entity.Name | ToLower}} body domain.{{.Entity.Name}} true "{{.Entity.Name}} object"
// @Success 201 {object} domain.{{.Entity.Name}}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /{{.Entity.Name | ToSnakeCase}}s [post]
func (c *{{.Entity.Name}}Controller) Create(ctx *gin.Context) {
	var entity domain.{{.Entity.Name}}
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.useCase.Create(ctx, &entity); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, entity)
}

// Get{{.Entity.Name}} godoc
// @Summary Get a {{.Entity.Name | ToLower}} by ID
// @Description Get a {{.Entity.Name | ToLower}} by its ID
// @Tags {{.Entity.Name | ToLower}}s
// @Accept json
// @Produce json
// @Param id path string true "{{.Entity.Name}} ID"
// @Success 200 {object} domain.{{.Entity.Name}}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /{{.Entity.Name | ToSnakeCase}}s/{id} [get]
func (c *{{.Entity.Name}}Controller) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	entity, err := c.useCase.Get(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if entity == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "{{.Entity.Name | ToLower}} not found"})
		return
	}

	ctx.JSON(http.StatusOK, entity)
}

// Update{{.Entity.Name}} godoc
// @Summary Update a {{.Entity.Name | ToLower}}
// @Description Update a {{.Entity.Name | ToLower}} with the input payload
// @Tags {{.Entity.Name | ToLower}}s
// @Accept json
// @Produce json
// @Param id path string true "{{.Entity.Name}} ID"
// @Param {{.Entity.Name | ToLower}} body domain.{{.Entity.Name}} true "{{.Entity.Name}} object"
// @Success 200 {object} domain.{{.Entity.Name}}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /{{.Entity.Name | ToSnakeCase}}s/{id} [put]
func (c *{{.Entity.Name}}Controller) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var entity domain.{{.Entity.Name}}
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity.ID = id
	if err := c.useCase.Update(ctx, &entity); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, entity)
}

// Delete{{.Entity.Name}} godoc
// @Summary Delete a {{.Entity.Name | ToLower}}
// @Description Delete a {{.Entity.Name | ToLower}} by its ID
// @Tags {{.Entity.Name | ToLower}}s
// @Accept json
// @Produce json
// @Param id path string true "{{.Entity.Name}} ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /{{.Entity.Name | ToSnakeCase}}s/{id} [delete]
func (c *{{.Entity.Name}}Controller) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	if err := c.useCase.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// List{{.Entity.Name}} godoc
// @Summary List all {{.Entity.Name | ToLower}}s
// @Description Get a list of all {{.Entity.Name | ToLower}}s
// @Tags {{.Entity.Name | ToLower}}s
// @Accept json
// @Produce json
// @Success 200 {array} domain.{{.Entity.Name}}
// @Failure 500 {object} map[string]interface{}
// @Router /{{.Entity.Name | ToSnakeCase}}s [get]
func (c *{{.Entity.Name}}Controller) List(ctx *gin.Context) {
	entities, err := c.useCase.List(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, entities)
}
`

	// Шаблоны для gRPC proto файла
	protoTemplate = `syntax = "proto3";

package {{.Entity.Name | ToLower}};

option go_package = "github.com/KulikovAR/{{.Module}}/pkg/proto/{{.Entity.Name | ToLower}}";

import "google/protobuf/timestamp.proto";

service {{.Entity.Name}}Service {
  rpc Create{{.Entity.Name}}(Create{{.Entity.Name}}Request) returns ({{.Entity.Name}}Response);
  rpc Get{{.Entity.Name}}(Get{{.Entity.Name}}Request) returns ({{.Entity.Name}}Response);
  rpc Update{{.Entity.Name}}(Update{{.Entity.Name}}Request) returns ({{.Entity.Name}}Response);
  rpc Delete{{.Entity.Name}}(Delete{{.Entity.Name}}Request) returns (Delete{{.Entity.Name}}Response);
  rpc List{{.Entity.Name}}s(List{{.Entity.Name}}sRequest) returns (List{{.Entity.Name}}sResponse);
}

message {{.Entity.Name}} {
  string id = 1;
  {{range $i, $field := .Entity.Fields}}
  {{$field.Type | ToProtoType}} {{$field.Name | ToSnakeCase}} = {{add $i 2}};
  {{end}}
  google.protobuf.Timestamp created_at = {{add (len .Entity.Fields) 2}};
  google.protobuf.Timestamp updated_at = {{add (len .Entity.Fields) 3}};
}

message Create{{.Entity.Name}}Request {
  {{range $i, $field := .Entity.Fields}}
  {{$field.Type | ToProtoType}} {{$field.Name | ToSnakeCase}} = {{add $i 1}};
  {{end}}
}

message Get{{.Entity.Name}}Request {
  string id = 1;
}

message Update{{.Entity.Name}}Request {
  string id = 1;
  {{range $i, $field := .Entity.Fields}}
  {{$field.Type | ToProtoType}} {{$field.Name | ToSnakeCase}} = {{add $i 2}};
  {{end}}
}

message Delete{{.Entity.Name}}Request {
  string id = 1;
}

message Delete{{.Entity.Name}}Response {
  bool success = 1;
}

message List{{.Entity.Name}}sRequest {
  int32 page = 1;
  int32 limit = 2;
}

message List{{.Entity.Name}}sResponse {
  repeated {{.Entity.Name}} {{.Entity.Name | ToLower}}s = 1;
  int32 total = 2;
}

message {{.Entity.Name}}Response {
  {{.Entity.Name}} {{.Entity.Name | ToLower}} = 1;
  string error = 2;
}
`

	// Шаблоны для gRPC controller
	grpcControllerTemplate = `package grpc

import (
	"context"
	"time"
	"github.com/KulikovAR/{{.Module}}/internal/domain"
	"github.com/KulikovAR/{{.Module}}/internal/usecase"
	"github.com/KulikovAR/{{.Module}}/pkg/proto/{{.Entity.Name | ToLower}}"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type {{.Entity.Name}}GRPCController struct {
	{{.Entity.Name | ToLower}}.Unimplemented{{.Entity.Name}}ServiceServer
	useCase usecase.{{.Entity.Name}}UseCase
}

func New{{.Entity.Name}}GRPCController(useCase usecase.{{.Entity.Name}}UseCase) *{{.Entity.Name}}GRPCController {
	return &{{.Entity.Name}}GRPCController{useCase: useCase}
}

func (c *{{.Entity.Name}}GRPCController) Create{{.Entity.Name}}(ctx context.Context, req *{{.Entity.Name | ToLower}}.Create{{.Entity.Name}}Request) (*{{.Entity.Name | ToLower}}.{{.Entity.Name}}Response, error) {
	entity := &domain.{{.Entity.Name}}{
		{{range .Entity.Fields}}
		{{.Name}}: req.Get{{.Name}}(),
		{{end}}
	}

	if err := c.useCase.Create(ctx, entity); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create {{.Entity.Name | ToLower}}: %v", err)
	}

	return &{{.Entity.Name | ToLower}}.{{.Entity.Name}}Response{
		{{.Entity.Name | ToLower}}: c.domainToProto(entity),
	}, nil
}

func (c *{{.Entity.Name}}GRPCController) Get{{.Entity.Name}}(ctx context.Context, req *{{.Entity.Name | ToLower}}.Get{{.Entity.Name}}Request) (*{{.Entity.Name | ToLower}}.{{.Entity.Name}}Response, error) {
	entity, err := c.useCase.Get(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get {{.Entity.Name | ToLower}}: %v", err)
	}

	if entity == nil {
		return nil, status.Errorf(codes.NotFound, "{{.Entity.Name | ToLower}} not found")
	}

	return &{{.Entity.Name | ToLower}}.{{.Entity.Name}}Response{
		{{.Entity.Name | ToLower}}: c.domainToProto(entity),
	}, nil
}

func (c *{{.Entity.Name}}GRPCController) Update{{.Entity.Name}}(ctx context.Context, req *{{.Entity.Name | ToLower}}.Update{{.Entity.Name}}Request) (*{{.Entity.Name | ToLower}}.{{.Entity.Name}}Response, error) {
	entity := &domain.{{.Entity.Name}}{
		ID: req.GetId(),
		{{range .Entity.Fields}}
		{{.Name}}: req.Get{{.Name}}(),
		{{end}}
	}

	if err := c.useCase.Update(ctx, entity); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update {{.Entity.Name | ToLower}}: %v", err)
	}

	return &{{.Entity.Name | ToLower}}.{{.Entity.Name}}Response{
		{{.Entity.Name | ToLower}}: c.domainToProto(entity),
	}, nil
}

func (c *{{.Entity.Name}}GRPCController) Delete{{.Entity.Name}}(ctx context.Context, req *{{.Entity.Name | ToLower}}.Delete{{.Entity.Name}}Request) (*{{.Entity.Name | ToLower}}.Delete{{.Entity.Name}}Response, error) {
	if err := c.useCase.Delete(ctx, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete {{.Entity.Name | ToLower}}: %v", err)
	}

	return &{{.Entity.Name | ToLower}}.Delete{{.Entity.Name}}Response{
		Success: true,
	}, nil
}

func (c *{{.Entity.Name}}GRPCController) List{{.Entity.Name}}s(ctx context.Context, req *{{.Entity.Name | ToLower}}.List{{.Entity.Name}}sRequest) (*{{.Entity.Name | ToLower}}.List{{.Entity.Name}}sResponse, error) {
	entities, err := c.useCase.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list {{.Entity.Name | ToLower}}s: %v", err)
	}

	var protoEntities []*{{.Entity.Name | ToLower}}.{{.Entity.Name}}
	for _, entity := range entities {
		protoEntities = append(protoEntities, c.domainToProto(entity))
	}

	return &{{.Entity.Name | ToLower}}.List{{.Entity.Name}}sResponse{
		{{.Entity.Name | ToLower}}s: protoEntities,
		Total: int32(len(protoEntities)),
	}, nil
}

func (c *{{.Entity.Name}}GRPCController) domainToProto(entity *domain.{{.Entity.Name}}) *{{.Entity.Name | ToLower}}.{{.Entity.Name}} {
	return &{{.Entity.Name | ToLower}}.{{.Entity.Name}}{
		Id: entity.ID,
		{{range .Entity.Fields}}
		{{.Name}}: entity.{{.Name}},
		{{end}}
		CreatedAt: timestamppb.New(entity.CreatedAt),
		UpdatedAt: timestamppb.New(entity.UpdatedAt),
	}
}
`

	// Шаблоны для тестов
	testTemplate = `package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/KulikovAR/{{.Module}}/internal/domain"
	"github.com/KulikovAR/{{.Module}}/internal/controller"
	"github.com/KulikovAR/{{.Module}}/internal/usecase"
)

type Mock{{.Entity.Name}}UseCase struct {
	mock.Mock
}

func (m *Mock{{.Entity.Name}}UseCase) Create(ctx context.Context, entity *domain.{{.Entity.Name}}) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *Mock{{.Entity.Name}}UseCase) Get(ctx context.Context, id string) (*domain.{{.Entity.Name}}, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.{{.Entity.Name}}), args.Error(1)
}

func (m *Mock{{.Entity.Name}}UseCase) Update(ctx context.Context, entity *domain.{{.Entity.Name}}) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *Mock{{.Entity.Name}}UseCase) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *Mock{{.Entity.Name}}UseCase) List(ctx context.Context) ([]*domain.{{.Entity.Name}}, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.{{.Entity.Name}}), args.Error(1)
}

func setup{{.Entity.Name}}Test() (*gin.Engine, *Mock{{.Entity.Name}}UseCase, *controller.{{.Entity.Name}}Controller) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	mockUseCase := &Mock{{.Entity.Name}}UseCase{}
	controller := controller.New{{.Entity.Name}}Controller(mockUseCase)
	
	api := router.Group("/api/v1")
	{
		api.POST("/{{.Entity.Name | ToSnakeCase}}s", controller.Create)
		api.GET("/{{.Entity.Name | ToSnakeCase}}s/:id", controller.Get)
		api.PUT("/{{.Entity.Name | ToSnakeCase}}s/:id", controller.Update)
		api.DELETE("/{{.Entity.Name | ToSnakeCase}}s/:id", controller.Delete)
		api.GET("/{{.Entity.Name | ToSnakeCase}}s", controller.List)
	}
	
	return router, mockUseCase, controller
}

func Test{{.Entity.Name}}Controller_Create(t *testing.T) {
	router, mockUseCase, _ := setup{{.Entity.Name}}Test()
	
	entity := &domain.{{.Entity.Name}}{
		{{range .Entity.Fields}}
		{{.Name}}: {{.Type | ToTestValue}},
		{{end}}
	}
	
	mockUseCase.On("Create", mock.Anything, mock.AnythingOfType("*domain.{{.Entity.Name}}"))
		.Return(nil)
	
	body, _ := json.Marshal(entity)
	req, _ := http.NewRequest("POST", "/api/v1/{{.Entity.Name | ToSnakeCase}}s", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	mockUseCase.AssertExpectations(t)
}

func Test{{.Entity.Name}}Controller_Get(t *testing.T) {
	router, mockUseCase, _ := setup{{.Entity.Name}}Test()
	
	entity := &domain.{{.Entity.Name}}{
		ID: "test-id",
		{{range .Entity.Fields}}
		{{.Name}}: {{.Type | ToTestValue}},
		{{end}}
	}
	
	mockUseCase.On("Get", mock.Anything, "test-id").Return(entity, nil)
	
	req, _ := http.NewRequest("GET", "/api/v1/{{.Entity.Name | ToSnakeCase}}s/test-id", nil)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func Test{{.Entity.Name}}Controller_Update(t *testing.T) {
	router, mockUseCase, _ := setup{{.Entity.Name}}Test()
	
	entity := &domain.{{.Entity.Name}}{
		ID: "test-id",
		{{range .Entity.Fields}}
		{{.Name}}: {{.Type | ToTestValue}},
		{{end}}
	}
	
	mockUseCase.On("Update", mock.Anything, mock.AnythingOfType("*domain.{{.Entity.Name}}"))
		.Return(nil)
	
	body, _ := json.Marshal(entity)
	req, _ := http.NewRequest("PUT", "/api/v1/{{.Entity.Name | ToSnakeCase}}s/test-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func Test{{.Entity.Name}}Controller_Delete(t *testing.T) {
	router, mockUseCase, _ := setup{{.Entity.Name}}Test()
	
	mockUseCase.On("Delete", mock.Anything, "test-id").Return(nil)
	
	req, _ := http.NewRequest("DELETE", "/api/v1/{{.Entity.Name | ToSnakeCase}}s/test-id", nil)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNoContent, w.Code)
	mockUseCase.AssertExpectations(t)
}

func Test{{.Entity.Name}}Controller_List(t *testing.T) {
	router, mockUseCase, _ := setup{{.Entity.Name}}Test()
	
	entities := []*domain.{{.Entity.Name}}{
		{
			ID: "test-id-1",
			{{range .Entity.Fields}}
			{{.Name}}: {{.Type | ToTestValue}},
			{{end}}
		},
		{
			ID: "test-id-2",
			{{range .Entity.Fields}}
			{{.Name}}: {{.Type | ToTestValue}},
			{{end}}
		},
	}
	
	mockUseCase.On("List", mock.Anything).Return(entities, nil)
	
	req, _ := http.NewRequest("GET", "/api/v1/{{.Entity.Name | ToSnakeCase}}s", nil)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}
`

	// Шаблоны для миграций
	postgresMigrationTemplate = `-- +migrate Up
CREATE TABLE {{.Entity.Name | ToSnakeCase}}s (
    id VARCHAR(36) PRIMARY KEY,
    {{- range .Entity.Fields}}
    {{.Name | ToSnakeCase}} {{.Type | ToPostgresType}}{{if .Required}} NOT NULL{{end}}{{if .Unique}} UNIQUE{{end}},
    {{- end}}
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_{{.Entity.Name | ToSnakeCase}}s_created_at ON {{.Entity.Name | ToSnakeCase}}s(created_at);

-- +migrate Down
DROP TABLE {{.Entity.Name | ToSnakeCase}}s;
`

	mongodbMigrationTemplate = `{
    "collection": "{{.Entity.Name | ToSnakeCase}}s",
    "indexes": [
        {
            "keys": {
                "id": 1
            },
            "options": {
                "unique": true
            }
        },
        {
            "keys": {
                "created_at": -1
            }
        }
    ]
}`
)
