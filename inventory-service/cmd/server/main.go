package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc" // стандартный gRPC

	"inventory-service/internal/config"
	inventorygrpc "inventory-service/internal/delivery/grpc" // алиас для вашего пакета
	"inventory-service/internal/kafka"
	"inventory-service/internal/repository/postgres"
	"inventory-service/internal/usecase"
	pb "inventory-service/pkg/protos/inventory"
)

func main() {
	cfg := config.LoadConfig()

	// Подключение к PostgreSQL
	connStr := "postgres://" + cfg.DBUser + ":" + cfg.DBPassword + "@" + cfg.DBHost + ":" + cfg.DBPort + "/" + cfg.DBName
	dbPool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Создание репозитория
	stockRepo := postgres.NewStockRepository(dbPool)

	// Создание usecase
	stockUsecase := usecase.NewStockUsecase(stockRepo)

	// Kafka consumer
	consumer, err := kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID, cfg.KafkaTopic, stockUsecase)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Printf("Kafka consumer error: %v", err)
		}
	}()

	// Создаём gRPC сервер (из стандартного пакета)
	grpcServer := grpc.NewServer()

	// Создаём ваш сервер (используем алиас inventorygrpc)
	inventoryServer := inventorygrpc.NewInventoryServer(stockUsecase)

	// Регистрируем сервис
	pb.RegisterInventoryServiceServer(grpcServer, inventoryServer)

	// Слушаем порт
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down...")
		cancel()
		grpcServer.GracefulStop()
	}()

	log.Printf("Inventory service listening on port %s", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
