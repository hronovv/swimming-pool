package main

import (
	_ "github.com/joho/godotenv/autoload"

	"context"
	"log"
	"os"
	"pool/internal/group"
	"pool/internal/lane"
	"pool/internal/payment"
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

	laneRepository := lane.NewRepository(db)
	laneService := lane.NewService(laneRepository)
	laneHandler := lane.NewHandler(laneService)

	paymentRepository := payment.NewRepository(db)
	paymentService := payment.NewService(paymentRepository)
	paymentHandler := payment.NewHandler(paymentService)

	groupRepository := group.NewRepository(db)
	groupService := group.NewService(groupRepository)
	groupHandler := group.NewHandler(groupService)

	port := ":8080"
	srv := server.New(port,
		userHandler,
		tariffHandler,
		laneHandler,
		paymentHandler,
		groupHandler,
	)

	log.Printf("Starting server on http://localhost%s", port)
	if err := srv.Run(); err != nil {
		log.Fatalf("Error running server %v", err)
	}
}
