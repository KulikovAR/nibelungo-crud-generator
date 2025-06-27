package usecase

const (
	// Шаблоны для Docker
	dockerfileTemplate = "FROM golang:1.24-alpine AS builder\n\n" +
		"WORKDIR /app\n\n" +
		"# Установка зависимостей\n" +
		"RUN apk add --no-cache git\n\n" +
		"# Копирование go mod файлов\n" +
		"COPY go.mod go.sum ./\n" +
		"RUN go mod download\n\n" +
		"# Копирование исходного кода\n" +
		"COPY . .\n\n" +
		"# Сборка приложения\n" +
		"RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server\n\n" +
		"# Финальный образ\n" +
		"FROM alpine:latest\n\n" +
		"RUN apk --no-cache add ca-certificates\n\n" +
		"WORKDIR /root/\n\n" +
		"# Копирование бинарного файла\n" +
		"COPY --from=builder /app/main .\n\n" +
		"# Копирование миграций\n" +
		"COPY --from=builder /app/migrations ./migrations\n\n" +
		"# Открытие порта\n" +
		"EXPOSE {{.Port}}\n\n" +
		"# Запуск приложения\n" +
		"CMD [\"./main\"]\n"

	dockerComposeTemplate = "version: '3.8'\n\n" +
		"services:\n" +
		"  app:\n" +
		"    build: .\n" +
		"    ports:\n" +
		"      - \"{{.Port}}:{{.Port}}\"\n" +
		"    environment:\n" +
		"      - DB_HOST=postgres\n" +
		"      - DB_PORT=5432\n" +
		"      - DB_USER=postgres\n" +
		"      - DB_PASSWORD=password\n" +
		"      - DB_NAME={{.Name | ToLower}}\n" +
		"      - GRPC_PORT=50051\n" +
		"    depends_on:\n" +
		"      - postgres\n" +
		"    networks:\n" +
		"      - {{.Name | ToLower}}-network\n\n" +
		"  postgres:\n" +
		"    image: postgres:15-alpine\n" +
		"    environment:\n" +
		"      - POSTGRES_USER=postgres\n" +
		"      - POSTGRES_PASSWORD=password\n" +
		"      - POSTGRES_DB={{.Name | ToLower}}\n" +
		"    ports:\n" +
		"      - \"5432:5432\"\n" +
		"    volumes:\n" +
		"      - postgres_data:/var/lib/postgresql/data\n" +
		"      - ./migrations/postgres:/docker-entrypoint-initdb.d\n" +
		"    networks:\n" +
		"      - {{.Name | ToLower}}-network\n\n" +
		"  mongodb:\n" +
		"    image: mongo:7\n" +
		"    environment:\n" +
		"      - MONGO_INITDB_ROOT_USERNAME=admin\n" +
		"      - MONGO_INITDB_ROOT_PASSWORD=password\n" +
		"    ports:\n" +
		"      - \"27017:27017\"\n" +
		"    volumes:\n" +
		"      - mongodb_data:/data/db\n" +
		"    networks:\n" +
		"      - {{.Name | ToLower}}-network\n\n" +
		"volumes:\n" +
		"  postgres_data:\n" +
		"  mongodb_data:\n\n" +
		"networks:\n" +
		"  {{.Name | ToLower}}-network:\n" +
		"    driver: bridge\n"

	// Шаблоны для main.go
	mainTemplate = "package main\n\n" +
		"import (\n" +
		"	\"context\"\n" +
		"	\"fmt\"\n" +
		"	\"log\"\n" +
		"	\"net/http\"\n" +
		"	\"os\"\n" +
		"	\"os/signal\"\n" +
		"	\"syscall\"\n" +
		"	\"time\"\n\n" +
		"	\"github.com/gin-gonic/gin\"\n" +
		"	\"github.com/spf13/viper\"\n" +
		"	\"google.golang.org/grpc\"\n" +
		"	\"google.golang.org/grpc/reflection\"\n" +
		"	\"github.com/KulikovAR/{{.Module}}/internal/controller\"\n" +
		"	\"github.com/KulikovAR/{{.Module}}/internal/repository\"\n" +
		"	\"github.com/KulikovAR/{{.Module}}/internal/usecase\"\n" +
		"	\"github.com/KulikovAR/{{.Module}}/pkg/database\"\n" +
		"	\"github.com/KulikovAR/{{.Module}}/pkg/grpc/server\"\n" +
		"	\"github.com/KulikovAR/{{.Module}}/pkg/grpc/controller/grpc\"\n" +
		"	\"github.com/KulikovAR/{{.Module}}/pkg/proto/{{range .Entities}}{{.Name | ToLower}}{{end}}\"\n" +
		")\n\n" +
		"func main() {\n" +
		"	// Загрузка конфигурации\n" +
		"	viper.SetConfigName(\"config\")\n" +
		"	viper.SetConfigType(\"yaml\")\n" +
		"	viper.AddConfigPath(\".\")\n" +
		"	viper.AutomaticEnv()\n\n" +
		"	if err := viper.ReadInConfig(); err != nil {\n" +
		"		log.Printf(\"Warning: Could not read config file: %v\", err)\n" +
		"	}\n\n" +
		"	// Установка значений по умолчанию\n" +
		"	viper.SetDefault(\"port\", {{.Port}})\n" +
		"	viper.SetDefault(\"db.host\", \"localhost\")\n" +
		"	viper.SetDefault(\"db.port\", 5432)\n" +
		"	viper.SetDefault(\"db.user\", \"postgres\")\n" +
		"	viper.SetDefault(\"db.password\", \"password\")\n" +
		"	viper.SetDefault(\"db.name\", \"{{.Name | ToLower}}\")\n\n" +
		"	// Подключение к базе данных\n" +
		"	db, err := database.Connect()\n" +
		"	if err != nil {\n" +
		"		log.Fatalf(\"Failed to connect to database: %v\", err)\n" +
		"	}\n" +
		"	defer db.Close()\n\n" +
		"	// Инициализация репозиториев\n" +
		"	{{range .Entities}}\n" +
		"	{{.Name | ToLower}}Repo := repository.New{{.Name}}Repository(db)\n" +
		"	{{end}}\n\n" +
		"	// Инициализация use cases\n" +
		"	{{range .Entities}}\n" +
		"	{{.Name | ToLower}}UseCase := usecase.New{{.Name}}UseCase({{.Name | ToLower}}Repo)\n" +
		"	{{end}}\n\n" +
		"	// Инициализация контроллеров\n" +
		"	{{range .Entities}}\n" +
		"	{{.Name | ToLower}}Controller := controller.New{{.Name}}Controller({{.Name | ToLower}}UseCase)\n" +
		"	{{end}}\n\n" +
		"	// Настройка Gin\n" +
		"	gin.SetMode(gin.ReleaseMode)\n" +
		"	router := gin.Default()\n\n" +
		"	// Middleware\n" +
		"	router.Use(gin.Logger())\n" +
		"	router.Use(gin.Recovery())\n\n" +
		"	// API маршруты\n" +
		"	api := router.Group(\"/api/v1\")\n" +
		"	{\n" +
		"		{{range .Entities}}\n" +
		"		{{.Name | ToLower}}s := api.Group(\"/{{.Name | ToSnakeCase}}s\")\n" +
		"		{\n" +
		"			{{.Name | ToLower}}s.POST(\"\", {{.Name | ToLower}}Controller.Create)\n" +
		"			{{.Name | ToLower}}s.GET(\"/:id\", {{.Name | ToLower}}Controller.Get)\n" +
		"			{{.Name | ToLower}}s.PUT(\"/:id\", {{.Name | ToLower}}Controller.Update)\n" +
		"			{{.Name | ToLower}}s.DELETE(\"/:id\", {{.Name | ToLower}}Controller.Delete)\n" +
		"			{{.Name | ToLower}}s.GET(\"\", {{.Name | ToLower}}Controller.List)\n" +
		"		}\n" +
		"		{{end}}\n" +
		"	}\n\n" +
		"	// Swagger документация\n" +
		"	if viper.GetBool(\"swagger\") {\n" +
		"		router.GET(\"/swagger/*any\", ginSwagger.WrapHandler(swaggerFiles.Handler))\n" +
		"	}\n\n" +
		"	// Запуск HTTP сервера\n" +
		"	port := viper.GetInt(\"port\")\n" +
		"	srv := &http.Server{\n" +
		"		Addr:    fmt.Sprintf(\":%d\", port),\n" +
		"		Handler: router,\n" +
		"	}\n\n" +
		"	go func() {\n" +
		"		log.Printf(\"Starting HTTP server on port %d\", port)\n" +
		"		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {\n" +
		"			log.Fatalf(\"Failed to start HTTP server: %v\", err)\n" +
		"		}\n" +
		"	}()\n\n" +
		"	{{if .Features.GRPC}}\n" +
		"	// Запуск gRPC сервера\n" +
		"	grpcPort := viper.GetInt(\"grpc.port\")\n" +
		"	grpcServer := grpc.NewServer()\n" +
		"	\n" +
		"	{{range .Entities}}\n" +
		"	{{.Name | ToLower}}GRPCController := grpc.New{{.Name}}GRPCController({{.Name | ToLower}}UseCase)\n" +
		"	{{.Name | ToLower}}.Register{{.Name}}ServiceServer(grpcServer, {{.Name | ToLower}}GRPCController)\n" +
		"	{{end}}\n" +
		"	\n" +
		"	reflection.Register(grpcServer)\n" +
		"	\n" +
		"	go func() {\n" +
		"		log.Printf(\"Starting gRPC server on port %d\", grpcPort)\n" +
		"		lis, err := net.Listen(\"tcp\", fmt.Sprintf(\":%d\", grpcPort))\n" +
		"		if err != nil {\n" +
		"			log.Fatalf(\"Failed to listen for gRPC: %v\", err)\n" +
		"		}\n" +
		"		if err := grpcServer.Serve(lis); err != nil {\n" +
		"			log.Fatalf(\"Failed to serve gRPC: %v\", err)\n" +
		"		}\n" +
		"	}()\n" +
		"	{{end}}\n\n" +
		"	// Graceful shutdown\n" +
		"	quit := make(chan os.Signal, 1)\n" +
		"	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)\n" +
		"	<-quit\n" +
		"	log.Println(\"Shutting down server...\")\n\n" +
		"	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)\n" +
		"	defer cancel()\n\n" +
		"	if err := srv.Shutdown(ctx); err != nil {\n" +
		"		log.Fatal(\"Server forced to shutdown:\", err)\n" +
		"	}\n\n" +
		"	{{if .Features.GRPC}}\n" +
		"	grpcServer.GracefulStop()\n" +
		"	{{end}}\n\n" +
		"	log.Println(\"Server exiting\")\n" +
		"}\n"

	// Шаблоны для конфигурации
	configYamlTemplate = "port: {{.Port}}\n" +
		"grpc:\n" +
		"  port: 50051\n\n" +
		"database:\n" +
		"  host: localhost\n" +
		"  port: 5432\n" +
		"  user: postgres\n" +
		"  password: password\n" +
		"  name: {{.Name | ToLower}}\n" +
		"  sslmode: disable\n\n" +
		"mongodb:\n" +
		"  uri: mongodb://admin:password@localhost:27017/{{.Name | ToLower}}?authSource=admin\n\n" +
		"swagger: true\n" +
		"log:\n" +
		"  level: info\n" +
		"  format: json\n"
)
