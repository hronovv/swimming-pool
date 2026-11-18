package main

import (
	_ "github.com/joho/godotenv/autoload"

	"context"
	"log"
	"os"
	"time"

	"pool/internal/auth"
	"pool/internal/client"
	"pool/internal/enrollment"
	"pool/internal/group"
	"pool/internal/journal"
	"pool/internal/lane"
	"pool/internal/payment"
	"pool/internal/permit"
	"pool/internal/postgres"
	"pool/internal/server"
	"pool/internal/service"
	"pool/internal/session"
	"pool/internal/standard"
	"pool/internal/tariff"
	"pool/internal/trainer"
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

	issuer, err := auth.NewIssuer(
		os.Getenv("JWT_SECRET"),
		15*time.Minute,
		7*24*time.Hour,
		os.Getenv("COOKIE_SECURE") == "true",
	)
	if err != nil {
		log.Fatalf("Error creating token issuer: %v", err)
	}
	authHandler := auth.NewHandler(userService, issuer)

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

	journalRepository := journal.NewRepository(db)
	journalService := journal.NewService(journalRepository)
	journalHandler := journal.NewHandler(journalService)

	serviceRepository := service.NewRepository(db)
	serviceService := service.NewService(serviceRepository)
	serviceHandler := service.NewHandler(serviceService)

	sessionRepository := session.NewRepository(db)
	sessionService := session.NewService(sessionRepository)
	sessionHandler := session.NewHandler(sessionService)

	clientRepository := client.NewRepository(db)
	clientService := client.NewService(clientRepository)
	clientHandler := client.NewHandler(clientService)

	trainerRepository := trainer.NewRepository(db)
	trainerService := trainer.NewService(trainerRepository)
	trainerHandler := trainer.NewHandler(trainerService)

	permitRepository := permit.NewRepository(db)
	permitService := permit.NewService(permitRepository)
	permitHandler := permit.NewHandler(permitService)

	enrollmentRepository := enrollment.NewRepository(db)
	enrollmentService := enrollment.NewService(enrollmentRepository)
	enrollmentHandler := enrollment.NewHandler(enrollmentService)

	standardRepository := standard.NewRepository(db)
	standardService := standard.NewService(standardRepository)
	standardHandler := standard.NewHandler(standardService)

	port := ":8080"
	srv := server.New(port,
		auth.RequireUser(userService, issuer),
		authHandler,
		userHandler,
		tariffHandler,
		laneHandler,
		paymentHandler,
		groupHandler,
		journalHandler,
		serviceHandler,
		sessionHandler,
		clientHandler,
		trainerHandler,
		permitHandler,
		enrollmentHandler,
		standardHandler,
	)

	log.Printf("Starting server on http://localhost%s", port)
	if err := srv.Run(); err != nil {
		log.Fatalf("Error running server %v", err)
	}
}
