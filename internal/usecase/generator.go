package usecase

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/KulikovAR/nibelungo-crud-generator/internal/domain"
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
	funcMap["ToProtoType"] = func(goType string) string {
		switch goType {
		case "string":
			return "string"
		case "int", "int32":
			return "int32"
		case "int64":
			return "int64"
		case "float32":
			return "float"
		case "float64":
			return "double"
		case "bool":
			return "bool"
		default:
			return "string"
		}
	}
	funcMap["ToTestValue"] = func(goType string) string {
		switch goType {
		case "string":
			return `"test-value"`
		case "int", "int32", "int64":
			return "123"
		case "float32", "float64":
			return "123.45"
		case "bool":
			return "true"
		default:
			return `"test-value"`
		}
	}
	funcMap["GetRandomPort"] = func() int {
		rand.Seed(time.Now().UnixNano())
		return rand.Intn(10000) + 8000
	}

	templates := make(map[string]*template.Template)

	// Базовые шаблоны
	templates["domain"] = template.Must(template.New("domain").Funcs(funcMap).Parse(domainEntityTemplate))
	templates["repository"] = template.Must(template.New("repository").Funcs(funcMap).Parse(repositoryInterfaceTemplate))
	templates["postgres"] = template.Must(template.New("postgres").Funcs(funcMap).Parse(postgresRepositoryTemplate))
	templates["mongodb"] = template.Must(template.New("mongodb").Funcs(funcMap).Parse(mongodbRepositoryTemplate))
	templates["usecase"] = template.Must(template.New("usecase").Funcs(funcMap).Parse(usecaseTemplate))
	templates["controller"] = template.Must(template.New("controller").Funcs(funcMap).Parse(restControllerTemplate))
	templates["postgres_migration"] = template.Must(template.New("postgres_migration").Funcs(funcMap).Parse(postgresMigrationTemplate))
	templates["mongodb_migration"] = template.Must(template.New("mongodb_migration").Funcs(funcMap).Parse(mongodbMigrationTemplate))

	// Новые шаблоны
	templates["rest_controller"] = template.Must(template.New("rest_controller").Funcs(funcMap).Parse(restControllerTemplate))
	templates["proto"] = template.Must(template.New("proto").Funcs(funcMap).Parse(protoTemplate))
	templates["test"] = template.Must(template.New("test").Funcs(funcMap).Parse(testTemplate))
	templates["dockerfile"] = template.Must(template.New("dockerfile").Funcs(funcMap).Parse(dockerfileTemplate))
	templates["docker_compose"] = template.Must(template.New("docker_compose").Funcs(funcMap).Parse(dockerComposeTemplate))
	templates["main"] = template.Must(template.New("main").Funcs(funcMap).Parse(mainTemplate))
	templates["config_yaml"] = template.Must(template.New("config_yaml").Funcs(funcMap).Parse(configYamlTemplate))

	return &generator{
		templates: templates,
	}
}

func (g *generator) Generate(config *domain.ProjectConfig) error {
	// Устанавливаем случайный порт если не указан
	if config.Port == 0 {
		config.Port = rand.Intn(10000) + 8000
	}

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

	// Генерируем общие файлы проекта
	if err := g.generateProjectFiles(config); err != nil {
		return fmt.Errorf("failed to generate project files: %w", err)
	}

	return nil
}

func (g *generator) createProjectStructure(config *domain.ProjectConfig) error {
	dirs := []string{
		"cmd/server",
		"internal/domain",
		"internal/repository",
		"internal/usecase",
		"internal/controller",
		"pkg/api",
		"pkg/database",
		"pkg/grpc/controller",
		"pkg/proto",
		"migrations/postgres",
		"migrations/mongodb",
		"tests",
		"docs",
	}

	// Добавляем директории для proto файлов
	if config.Features.GRPC {
		dirs = append(dirs, "proto")
	}

	// Добавляем директории для тестов
	if config.Features.Tests {
		dirs = append(dirs, "tests/unit", "tests/integration")
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(config.Name, dir), 0755); err != nil {
			return err
		}
	}

	// Создаем go.mod
	goModContent := fmt.Sprintf("module %s\n\ngo 1.24\n\nrequire (\n", config.Module)
	goModContent += "	github.com/gin-gonic/gin v1.9.1\n"
	goModContent += "	github.com/spf13/viper v1.18.2\n"
	goModContent += "	github.com/lib/pq v1.10.9\n"
	goModContent += "	go.mongodb.org/mongo-driver v1.13.1\n"
	goModContent += "	github.com/golang-migrate/migrate/v4 v4.17.0\n"
	goModContent += "	github.com/stretchr/testify v1.8.4\n"
	goModContent += "	github.com/google/uuid v1.6.0\n"

	if config.Features.GRPC {
		goModContent += "	google.golang.org/grpc v1.62.1\n"
		goModContent += "	google.golang.org/protobuf v1.33.0\n"
	}

	if config.Features.Swagger {
		goModContent += "	github.com/swaggo/gin-swagger v1.6.0\n"
		goModContent += "	github.com/swaggo/files v1.0.1\n"
		goModContent += "	github.com/swaggo/swag v1.16.3\n"
	}

	goModContent += ")\n"

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
	if err := g.generateFile(config, "repository", struct {
		Entity domain.Entity
		Module string
	}{entity, config.Module}, filepath.Join("internal/repository", strcase.ToSnake(entity.Name)+".go")); err != nil {
		return err
	}

	// Генерируем usecase
	if err := g.generateFile(config, "usecase", struct {
		Entity domain.Entity
		Module string
	}{entity, config.Module}, filepath.Join("internal/usecase", strcase.ToSnake(entity.Name)+".go")); err != nil {
		return err
	}

	// Генерируем REST controller
	if config.Features.REST {
		if err := g.generateFile(config, "rest_controller", struct {
			Entity domain.Entity
			Module string
		}{entity, config.Module}, filepath.Join("internal/controller", strcase.ToSnake(entity.Name)+".go")); err != nil {
			return err
		}
	}

	// Генерируем proto файлы
	if config.Features.GRPC {
		if err := g.generateFile(config, "proto", struct {
			Entity domain.Entity
			Module string
		}{entity, config.Module}, filepath.Join("proto", strcase.ToSnake(entity.Name)+".proto")); err != nil {
			return err
		}
	}

	// Генерируем тесты
	if config.Features.Tests {
		if err := g.generateFile(config, "test", struct {
			Entity domain.Entity
			Module string
		}{entity, config.Module}, filepath.Join("tests/unit", strcase.ToSnake(entity.Name)+"_controller_test.go")); err != nil {
			return err
		}
	}

	// Генерируем репозитории
	for _, repo := range config.Repositories {
		if repo == "postgres" || repo == "mongodb" {
			if err := g.generateFile(config, repo, struct {
				Entity domain.Entity
				Module string
			}{entity, config.Module}, filepath.Join("internal/repository", repo, strcase.ToSnake(entity.Name)+".go")); err != nil {
				return err
			}
		} else {
			if err := g.generateFile(config, repo, entity, filepath.Join("internal/repository", repo, strcase.ToSnake(entity.Name)+".go")); err != nil {
				return err
			}
		}
	}

	// Генерируем миграции
	if config.Features.Migrations {
		if err := g.generateFile(config, "postgres_migration", struct {
			Entity domain.Entity
			Module string
		}{entity, config.Module}, filepath.Join("migrations/postgres", fmt.Sprintf("001_create_%s.up.sql", strcase.ToSnake(entity.Name)))); err != nil {
			return err
		}

		if err := g.generateFile(config, "mongodb_migration", struct {
			Entity domain.Entity
			Module string
		}{entity, config.Module}, filepath.Join("migrations/mongodb", fmt.Sprintf("001_create_%s.up.json", strcase.ToSnake(entity.Name)))); err != nil {
			return err
		}
	}

	return nil
}

func (g *generator) generateProjectFiles(config *domain.ProjectConfig) error {
	// Генерируем main.go
	if err := g.generateFile(config, "main", config, filepath.Join("cmd/server", "main.go")); err != nil {
		return err
	}

	// Генерируем конфигурацию
	if err := g.generateFile(config, "config_yaml", config, "config.yaml"); err != nil {
		return err
	}

	// Генерируем Docker файлы
	if config.Features.Docker {
		if err := g.generateFile(config, "dockerfile", config, "Dockerfile"); err != nil {
			return err
		}

		if err := g.generateFile(config, "docker_compose", config, "docker-compose.yml"); err != nil {
			return err
		}
	}

	// Генерируем README
	if err := g.generateREADME(config); err != nil {
		return err
	}

	// Генерируем Makefile
	if err := g.generateMakefile(config); err != nil {
		return err
	}

	return nil
}

func (g *generator) generateFile(config *domain.ProjectConfig, templateName string, data interface{}, path string) error {
	tmpl, ok := g.templates[templateName]
	if !ok {
		return fmt.Errorf("template %s not found", templateName)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	fullPath := filepath.Join(config.Name, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(fullPath, buf.Bytes(), 0644)
}

func (g *generator) generateREADME(config *domain.ProjectConfig) error {
	readmeContent := fmt.Sprintf("# %s\n\n", config.Name)
	readmeContent += "Автоматически сгенерированный CRUD сервис на Go.\n\n"
	readmeContent += "## Возможности\n\n"
	readmeContent += "- REST API с автоматической документацией\n"

	if config.Features.GRPC {
		readmeContent += "- gRPC API\n"
	}

	readmeContent += "- Поддержка PostgreSQL и MongoDB\n"
	readmeContent += "- Автоматические тесты\n"
	readmeContent += "- Docker окружение\n"
	readmeContent += "- Миграции базы данных\n\n"
	readmeContent += "## Быстрый старт\n\n"
	readmeContent += "### С Docker\n\n"
	readmeContent += "```bash\n"
	readmeContent += "docker-compose up -d\n"
	readmeContent += "```\n\n"
	readmeContent += fmt.Sprintf("Сервис будет доступен на порту %d\n\n", config.Port)
	readmeContent += "### Без Docker\n\n"
	readmeContent += "1. Установите зависимости:\n"
	readmeContent += "```bash\n"
	readmeContent += "go mod download\n"
	readmeContent += "```\n\n"
	readmeContent += "2. Настройте базу данных\n\n"
	readmeContent += "3. Запустите сервис:\n"
	readmeContent += "```bash\n"
	readmeContent += "go run cmd/server/main.go\n"
	readmeContent += "```\n\n"
	readmeContent += "## API Endpoints\n\n"

	for _, entity := range config.Entities {
		readmeContent += fmt.Sprintf("### %s\n\n", entity.Name)
		readmeContent += fmt.Sprintf("- POST /api/v1/%ss - Создать %s\n", strcase.ToSnake(entity.Name), entity.Name)
		readmeContent += fmt.Sprintf("- GET /api/v1/%ss/:id - Получить %s по ID\n", strcase.ToSnake(entity.Name), entity.Name)
		readmeContent += fmt.Sprintf("- PUT /api/v1/%ss/:id - Обновить %s\n", strcase.ToSnake(entity.Name), entity.Name)
		readmeContent += fmt.Sprintf("- DELETE /api/v1/%ss/:id - Удалить %s\n", strcase.ToSnake(entity.Name), entity.Name)
		readmeContent += fmt.Sprintf("- GET /api/v1/%ss - Список всех %s\n\n", strcase.ToSnake(entity.Name), entity.Name)
	}

	readmeContent += "## Тестирование\n\n"
	readmeContent += "```bash\n"
	readmeContent += "go test ./tests/...\n"
	readmeContent += "```\n\n"
	readmeContent += "## Документация API\n\n"
	readmeContent += fmt.Sprintf("Swagger UI доступен по адресу: http://localhost:%d/swagger/index.html\n", config.Port)

	return os.WriteFile(filepath.Join(config.Name, "README.md"), []byte(readmeContent), 0644)
}

func (g *generator) generateMakefile(config *domain.ProjectConfig) error {
	makefileContent := fmt.Sprintf("# Makefile для %s\n\n", config.Name)
	makefileContent += ".PHONY: build run test clean docker-build docker-run\n\n"
	makefileContent += "# Сборка приложения\n"
	makefileContent += "build:\n"
	makefileContent += "	go build -o bin/server cmd/server/main.go\n\n"
	makefileContent += "# Запуск приложения\n"
	makefileContent += "run:\n"
	makefileContent += "	go run cmd/server/main.go\n\n"
	makefileContent += "# Запуск тестов\n"
	makefileContent += "test:\n"
	makefileContent += "	go test -v ./tests/...\n\n"
	makefileContent += "# Очистка\n"
	makefileContent += "clean:\n"
	makefileContent += "	rm -rf bin/\n"
	makefileContent += "	go clean\n\n"
	makefileContent += "# Сборка Docker образа\n"
	makefileContent += "docker-build:\n"
	makefileContent += fmt.Sprintf("	docker build -t %s .\n\n", config.Name)
	makefileContent += "# Запуск в Docker\n"
	makefileContent += "docker-run:\n"
	makefileContent += "	docker-compose up -d\n\n"
	makefileContent += "# Остановка Docker\n"
	makefileContent += "docker-stop:\n"
	makefileContent += "	docker-compose down\n\n"
	makefileContent += "# Просмотр логов\n"
	makefileContent += "logs:\n"
	makefileContent += "	docker-compose logs -f\n\n"
	makefileContent += "# Миграции\n"
	makefileContent += "migrate-up:\n"
	makefileContent += fmt.Sprintf("	migrate -path migrations/postgres -database \"postgres://postgres:password@localhost:5432/%s?sslmode=disable\" up\n\n", config.Name)
	makefileContent += "migrate-down:\n"
	makefileContent += fmt.Sprintf("	migrate -path migrations/postgres -database \"postgres://postgres:password@localhost:5432/%s?sslmode=disable\" down\n", config.Name)

	return os.WriteFile(filepath.Join(config.Name, "Makefile"), []byte(makefileContent), 0644)
}
