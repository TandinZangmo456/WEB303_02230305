package main

import (
	"log"
	"net"
	"os"

	"order-service/database"
	ordergrpc "order-service/grpc"
	menuv1 "order-service/proto/menuv1"
	orderv1 "order-service/proto/orderv1"
	userv1 "order-service/proto/userv1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Get database connection string from environment
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "orderdb")

	dsn := "host=" + dbHost + " port=" + dbPort + " user=" + dbUser + 
		" password=" + dbPassword + " dbname=" + dbName + " sslmode=disable"

	// Connect to database
	if err := database.Connect(dsn); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Connect to user service
	userServiceAddr := getEnv("USER_SERVICE_ADDR", "localhost:50051")
	userConn, err := grpc.Dial(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	defer userConn.Close()
	userClient := userv1.NewUserServiceClient(userConn)

	// Connect to menu service
	menuServiceAddr := getEnv("MENU_SERVICE_ADDR", "localhost:50052")
	menuConn, err := grpc.Dial(menuServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to menu service: %v", err)
	}
	defer menuConn.Close()
	menuClient := menuv1.NewMenuServiceClient(menuConn)

	// Create gRPC server
	grpcPort := getEnv("GRPC_PORT", "50053")
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	orderv1.RegisterOrderServiceServer(s, ordergrpc.NewOrderServer(userClient, menuClient))

	log.Printf("Order service listening on port %s", grpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
