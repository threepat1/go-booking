package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	customerhttp "booking-app/modules/customer/interface/customerhttp"
	customerRepository "booking-app/modules/customer/repository"
	customerUsecase "booking-app/modules/customer/usecase"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func main() {
	// Set up MongoDB repository
	// Load environment variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Set up MongoDB repository
	repo, err := customerRepository.NewMongoCustomerRepository("mongodb://localhost:27017", "customerdb", "customers")
	if err != nil {
		log.Fatal(err)
	}

	// Read SMTP configuration from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatalf("Invalid SMTP port: %v", err)
	}
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	// Set up email sender
	emailSender := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPassword)

	// Set up use case
	appURL := "http://localhost:8080" // or read from environment variable
	cu := customerUsecase.NewCustomerUsecase(repo, emailSender, appURL)

	// Set up HTTP handler
	ch := customerhttp.NewCustomerHandler(cu)

	// Set up router
	r := mux.NewRouter()
	customerhttp.RegisterCustomerHandlers(r, ch)

	// Start server
	log.Fatal(http.ListenAndServe(":8080", r))
}
