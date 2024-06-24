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
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := ch.CustomerUsecase.CreateCustomer(context.Background(), &customer); err != nil {
		if err == usecase.ErrEmailAlreadyExists {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	json.NewEncoder(w).Encode(customer)
}

func (ch *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var customer domain.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := ch.CustomerUsecase.UpdateCustomer(context.Background(), params["id"], &customer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(customer)
}

func RegisterCustomerHandlers(router *mux.Router, ch *CustomerHandler) {
	router.HandleFunc("/customers", ch.GetCustomers).Methods("GET")
	router.HandleFunc("/customers", ch.CreateCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", ch.UpdateCustomer).Methods("PUT")
}
