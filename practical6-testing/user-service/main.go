package main

import (
	"log"
	"net"
	"os"

	"user-service/database"
	usergrpc "user-service/grpc"
	userv1 "user-service/proto/userv1"

	"google.golang.org/grpc"
)

func main() {
	// Get database connection string from environment
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "userdb")

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

	// Create gRPC server
	grpcPort := getEnv("GRPC_PORT", "50051")
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	userv1.RegisterUserServiceServer(s, usergrpc.NewUserServer())

	log.Printf("User service listening on port %s", grpcPort)
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
