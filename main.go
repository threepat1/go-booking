package main

import (
	"log"
	"net/http"

	customerhttp "booking-app/modules/customer/interface/customerhttp"
	customerRepository "booking-app/modules/customer/repository"
	customerUsecase "booking-app/modules/customer/usecase"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Set up MongoDB repository
	repo, err := customerRepository.NewMongoCustomerRepository("mongodb://localhost:27017", "customerdb", "customers")
	if err != nil {
		log.Fatal(err)
	}

	cu := customerUsecase.NewCustomerUsecase(repo)

	// Set up HTTP handler
	ch := customerhttp.NewCustomerHandler(cu)

	// Set up router
	r := mux.NewRouter()
	customerhttp.RegisterCustomerHandlers(r, ch)

	// Start server
	log.Fatal(http.ListenAndServe(":8080", r))
}
