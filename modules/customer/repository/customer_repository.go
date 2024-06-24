package repository

import (
	"context"

	"booking-app/modules/customer/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerRepository interface {
	FindAll(ctx context.Context) ([]domain.Customer, error)
	Save(ctx context.Context, customer *domain.Customer) error
	Update(ctx context.Context, id primitive.ObjectID, customer *domain.Customer) error
	FindByEmail(ctx context.Context, email string) (*domain.Customer, error)
}
