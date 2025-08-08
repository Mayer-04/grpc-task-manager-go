package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Mayer-04/grpc-task-manager-go/internal/tasks/application"
	"github.com/Mayer-04/grpc-task-manager-go/internal/tasks/infrastructure"
	"github.com/Mayer-04/grpc-task-manager-go/pkg/taskpb"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Configuración del puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	// Configuración de la base de datos
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DB")

	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbName == "" {
		dbName = "taskdb"
	}

	// Crear connection string para PostgreSQL
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Conectar a PostgreSQL
	dbPool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Verificar conexión
	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to PostgreSQL")

	// Inicializar capas
	taskRepo := infrastructure.NewTaskRepository(dbPool)
	taskService := application.NewTaskService(taskRepo)
	taskHandler := infrastructure.NewTaskHandler(taskService)

	// Configurar servidor gRPC
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()

	// Registrar servicios
	taskpb.RegisterTaskServiceServer(grpcServer, taskHandler)

	// Habilitar reflection para herramientas como grpcui
	reflection.Register(grpcServer)

	// Canal para manejar shutdown graceful
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor en una goroutine
	go func() {
		log.Printf("gRPC server starting on port %s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	// Esperar señal de shutdown
	<-quit
	log.Println("Shutting down gRPC server...")

	// Graceful shutdown
	grpcServer.GracefulStop()
	log.Println("gRPC server stopped")
}
