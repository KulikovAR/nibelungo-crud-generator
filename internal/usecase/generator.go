package usecase

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/KulikovAR/go-core/internal/domain"
	"github.com/Masterminds/sprig/v3"
	"github.com/iancoleman/strcase"
)

type Generator interface {
	Generate(config *domain.ProjectConfig) error
}

type generator struct {
	templates map[string]*template.Template
}

func NewGenerator() Generator {
	funcMap := sprig.FuncMap()
	funcMap["ToLower"] = strings.ToLower
	funcMap["ToSnakeCase"] = strcase.ToSnake
	funcMap["ToPostgresType"] = func(goType string) string {
		switch goType {
		case "string":
			return "VARCHAR(255)"
		case "int", "int32", "int64":
			return "INTEGER"
		case "float32", "float64":
			return "FLOAT"
		case "bool":
			return "BOOLEAN"
		default:
			return "TEXT"
		}
	}

	templates := make(map[string]*template.Template)
	templates["domain"] = template.Must(template.New("domain").Funcs(funcMap).Parse(domainEntityTemplate))
	templates["repository"] = template.Must(template.New("repository").Funcs(funcMap).Parse(repositoryInterfaceTemplate))
	templates["postgres"] = template.Must(template.New("postgres").Funcs(funcMap).Parse(postgresRepositoryTemplate))
	templates["mongodb"] = template.Must(template.New("mongodb").Funcs(funcMap).Parse(mongodbRepositoryTemplate))
	templates["usecase"] = template.Must(template.New("usecase").Funcs(funcMap).Parse(usecaseTemplate))
	templates["controller"] = template.Must(template.New("controller").Funcs(funcMap).Parse(controllerTemplate))
	templates["postgres_migration"] = template.Must(template.New("postgres_migration").Funcs(funcMap).Parse(postgresMigrationTemplate))
	templates["mongodb_migration"] = template.Must(template.New("mongodb_migration").Funcs(funcMap).Parse(mongodbMigrationTemplate))

	return &generator{
		templates: templates,
	}
}

func (g *generator) Generate(config *domain.ProjectConfig) error {
	// Создаем структуру проекта
	if err := g.createProjectStructure(config); err != nil {
		return fmt.Errorf("failed to create project structure: %w", err)
	}

	// Генерируем код для каждой сущности
	for _, entity := range config.Entities {
		if err := g.generateEntityCode(config, entity); err != nil {
			return fmt.Errorf("failed to generate code for entity %s: %w", entity.Name, err)
		}
	}

	return nil
}

func (g *generator) createProjectStructure(config *domain.ProjectConfig) error {
	dirs := []string{
		"cmd",
		"internal/domain",
		"internal/repository",
		"internal/usecase",
		"internal/controller",
		"pkg/api",
		"migrations/postgres",
		"migrations/mongodb",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(config.Name, dir), 0755); err != nil {
			return err
		}
	}

	// Создаем go.mod
	goModContent := fmt.Sprintf(`module %s

go 1.24

require (
	github.com/spf13/cobra v1.8.0
	github.com/spf13/viper v1.18.2
	github.com/lib/pq v1.10.9
	go.mongodb.org/mongo-driver v1.13.1
)`, config.Module)

	if err := os.WriteFile(filepath.Join(config.Name, "go.mod"), []byte(goModContent), 0644); err != nil {
		return err
	}

	return nil
}

func (g *generator) generateEntityCode(config *domain.ProjectConfig, entity domain.Entity) error {
	// Генерируем domain
	if err := g.generateFile(config, "domain", entity, filepath.Join("internal/domain", strcase.ToSnake(entity.Name)+".go")); err != nil {
		return err
	}

	// Генерируем repository interface
	if err := g.generateFile(config, "repository", entity, filepath.Join("internal/repository", strcase.ToSnake(entity.Name)+".go")); err != nil {
		return err
	}

	// Генерируем usecase
	if err := g.generateFile(config, "usecase", entity, filepath.Join("internal/usecase", strcase.ToSnake(entity.Name)+".go")); err != nil {
		return err
	}

	// Генерируем controller
	if err := g.generateFile(config, "controller", entity, filepath.Join("internal/controller", strcase.ToSnake(entity.Name)+".go")); err != nil {
		return err
	}

	// Генерируем репозитории
	for _, repo := range config.Repositories {
		if err := g.generateFile(config, repo, entity, filepath.Join("internal/repository", repo, strcase.ToSnake(entity.Name)+".go")); err != nil {
			return err
		}
	}

	// Генерируем миграции
	if err := g.generateFile(config, "postgres_migration", entity, filepath.Join("migrations/postgres", fmt.Sprintf("001_create_%s.up.sql", strcase.ToSnake(entity.Name)))); err != nil {
		return err
	}

	if err := g.generateFile(config, "mongodb_migration", entity, filepath.Join("migrations/mongodb", fmt.Sprintf("001_create_%s.up.json", strcase.ToSnake(entity.Name)))); err != nil {
		return err
	}

	return nil
}

func (g *generator) generateFile(config *domain.ProjectConfig, templateName string, entity domain.Entity, path string) error {
	tmpl, ok := g.templates[templateName]
	if !ok {
		return fmt.Errorf("template %s not found", templateName)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, entity); err != nil {
		return err
	}

	fullPath := filepath.Join(config.Name, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(fullPath, buf.Bytes(), 0644)
}
