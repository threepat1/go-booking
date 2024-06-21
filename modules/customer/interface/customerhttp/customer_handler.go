package customerhttp

import (
	"context"
	"encoding/json"
	"net/http"

	"booking-app/modules/customer/domain"
	"booking-app/modules/customer/usecase"

	"github.com/gorilla/mux"
)

type CustomerHandler struct {
	CustomerUsecase *usecase.CustomerUsecase
}

func NewCustomerHandler(cu *usecase.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{CustomerUsecase: cu}
}

func (ch *CustomerHandler) GetCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := ch.CustomerUsecase.GetCustomers(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(customers)
}

func (ch *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer domain.Customer
	_ = json.NewDecoder(r.Body).Decode(&customer)
	if err := ch.CustomerUsecase.CreateCustomer(context.Background(), &customer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(customer)
}

func (ch *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var customer domain.Customer
	_ = json.NewDecoder(r.Body).Decode(&customer)
	if err := ch.CustomerUsecase.UpdateCustomer(context.Background(), params["id"], &customer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(customer)
}

func (ch *CustomerHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}
	if err := ch.CustomerUsecase.VerifyEmail(context.Background(), token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email verified successfully"))
}

func RegisterCustomerHandlers(router *mux.Router, ch *CustomerHandler) {
	router.HandleFunc("/customers", ch.GetCustomers).Methods("GET")
	router.HandleFunc("/customers", ch.CreateCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", ch.UpdateCustomer).Methods("PUT")
	router.HandleFunc("/verify-email", ch.VerifyEmail).Methods("GET")
}
