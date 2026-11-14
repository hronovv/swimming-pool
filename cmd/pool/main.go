package main

import (
	_ "github.com/joho/godotenv/autoload"

	"context"
	"log"
	"os"
	"pool/internal/postgres"
	"pool/internal/server"
	"pool/internal/tariff"
	"pool/internal/user"
)

func main() {

	ctx := context.Background()
	dsn := os.Getenv("DSN")

	db, err := postgres.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
	}
	log.Println("Connected to db pool!")
	defer db.Close()

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	tariffRepository := tariff.NewRepository(db)
	tariffService := tariff.NewService(tariffRepository)
	tariffHandler := tariff.NewHandler(tariffService)

	port := ":8080"
	srv := server.New(port, userHandler, tariffHandler)
	log.Printf("Starting server on %s", port)
	if err := srv.Run(); err != nil {
		log.Fatalf("Error running server %v", err)
	}
}
