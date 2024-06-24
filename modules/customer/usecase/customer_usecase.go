package usecase

import (
	"context"
	"errors"

	"time"

	"booking-app/modules/customer/domain"
	"booking-app/modules/customer/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type CustomerUsecase struct {
	CustomerRepo repository.CustomerRepository
}

func NewCustomerUsecase(cr repository.CustomerRepository) *CustomerUsecase {
	return &CustomerUsecase{CustomerRepo: cr}
}

func (cu *CustomerUsecase) GetCustomers(ctx context.Context) ([]domain.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return cu.CustomerRepo.FindAll(ctx)
}

func (cu *CustomerUsecase) CreateCustomer(ctx context.Context, customer *domain.Customer) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	customer.Password = string(hashedPassword)
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()

	// Save the customer
	customer.ID = primitive.NewObjectID()
	err = cu.CustomerRepo.Save(ctx, customer)
	if err != nil {
		if err == repository.ErrDuplicateKey {
			return ErrEmailAlreadyExists
		}
		return err
	}
	return nil
}

func (cu *CustomerUsecase) UpdateCustomer(ctx context.Context, id string, customer *domain.Customer) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	customer.ID = oid
	customer.UpdatedAt = time.Now()
	return cu.CustomerRepo.Update(ctx, oid, customer)
}
