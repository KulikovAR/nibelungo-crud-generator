package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/internal/controller"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/internal/repository"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/internal/usecase"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/pkg/database"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/pkg/grpc/server"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/pkg/grpc/controller/grpc"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/pkg/proto/userproduct"
)

func main() {
	// Загрузка конфигурации
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file: %v", err)
	}

	// Установка значений по умолчанию
	viper.SetDefault("port", 8080)
	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", 5432)
	viper.SetDefault("db.user", "postgres")
	viper.SetDefault("db.password", "password")
	viper.SetDefault("db.name", "some-project")

	// Подключение к базе данных
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Инициализация репозиториев
	
	userRepo := repository.NewuserRepository(db)
	
	productRepo := repository.NewproductRepository(db)
	

	// Инициализация use cases
	
	userUseCase := usecase.NewuserUseCase(userRepo)
	
	productUseCase := usecase.NewproductUseCase(productRepo)
	

	// Инициализация контроллеров
	
	userController := controller.NewuserController(userUseCase)
	
	productController := controller.NewproductController(productUseCase)
	

	// Настройка Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// API маршруты
	api := router.Group("/api/v1")
	{
		
		users := api.Group("/users")
		{
			users.POST("", userController.Create)
			users.GET("/:id", userController.Get)
			users.PUT("/:id", userController.Update)
			users.DELETE("/:id", userController.Delete)
			users.GET("", userController.List)
		}
		
		products := api.Group("/products")
		{
			products.POST("", productController.Create)
			products.GET("/:id", productController.Get)
			products.PUT("/:id", productController.Update)
			products.DELETE("/:id", productController.Delete)
			products.GET("", productController.List)
		}
		
	}

	// Swagger документация
	if viper.GetBool("swagger") {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Запуск HTTP сервера
	port := viper.GetInt("port")
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	go func() {
		log.Printf("Starting HTTP server on port %d", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	
	// Запуск gRPC сервера
	grpcPort := viper.GetInt("grpc.port")
	grpcServer := grpc.NewServer()
	
	
	userGRPCController := grpc.NewuserGRPCController(userUseCase)
	user.RegisteruserServiceServer(grpcServer, userGRPCController)
	
	productGRPCController := grpc.NewproductGRPCController(productUseCase)
	product.RegisterproductServiceServer(grpcServer, productGRPCController)
	
	
	reflection.Register(grpcServer)
	
	go func() {
		log.Printf("Starting gRPC server on port %d", grpcPort)
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
		if err != nil {
			log.Fatalf("Failed to listen for gRPC: %v", err)
		}
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()
	

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	
	grpcServer.GracefulStop()
	

	log.Println("Server exiting")
}
