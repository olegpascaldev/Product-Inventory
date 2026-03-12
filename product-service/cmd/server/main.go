package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"product-service/internal/clients/inventory"
	"product-service/internal/config"
	grpcproduct "product-service/internal/delivery/grpc"
	"product-service/internal/kafka"
	"product-service/internal/repository/postgres"
	"product-service/internal/usecase"
	pb "product-service/pkg/protos/product"
	"syscall"

	"google.golang.org/grpc"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	cfg := config.LoadConfig()
	// Connect to PostgreSQL
	connStr := "postgres://" + cfg.DBUser + ":" + cfg.DBPassword + "@" + cfg.DBHost + ":" + cfg.DBPort + "/" + cfg.DBName
	dbPool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Create repository
	productRepo := postgres.NewProductRepository(dbPool)

	// Create Kafka producer
	kafkaProducer, err := kafka.NewProducer(cfg.KafkaBrokers, cfg.KafkaTopic)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	// Create gRPC client for Inventory Service
	inventoryClient, err := inventory.NewClient(cfg.InventoryServiceAddr)
	if err != nil {
		log.Fatalf("Failed to create inventory client: %v", err)
	}
	defer inventoryClient.Close()

	// Create usecase
	productUsecase := usecase.NewProductUsecase(productRepo, kafkaProducer, inventoryClient)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	productServer := grpcproduct.NewProductServer(productUsecase)
	pb.RegisterProductServiceServer(grpcServer, productServer)

	// Listen
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	log.Printf("Product service listening on port %s", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
