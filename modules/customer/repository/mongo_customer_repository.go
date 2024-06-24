package repository

import (
	"context"
	"errors"
	"time"

	"booking-app/modules/customer/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrDuplicateKey = errors.New("duplicate key error")
)

type MongoCustomerRepository struct {
	collection *mongo.Collection
}

func NewMongoCustomerRepository(uri, dbName, collectionName string) (*MongoCustomerRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)
	return &MongoCustomerRepository{collection: collection}, nil
}

func (m *MongoCustomerRepository) FindAll(ctx context.Context) ([]domain.Customer, error) {
	cursor, err := m.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var customers []domain.Customer
	if err = cursor.All(ctx, &customers); err != nil {
		return nil, err
	}
	return customers, nil
}

func (m *MongoCustomerRepository) FindByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	var customer domain.Customer
	err := m.collection.FindOne(ctx, bson.M{"email": email}).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Email not found, no error
		}
		return nil, err // Other error occurred
	}
	return &customer, nil // Email found
}
func (m *MongoCustomerRepository) Save(ctx context.Context, customer *domain.Customer) error {
	// Check if email already exists
	existingCustomer, err := m.FindByEmail(ctx, customer.Email)
	if err != nil {
		return err // Return error if FindByEmail fails
	}
	if existingCustomer != nil {
		return mongo.ErrMultipleIndexDrop // Return conflict error if email already exists
	}

	// Insert new customer
	_, err = m.collection.InsertOne(ctx, customer)
	return err
}

func (m *MongoCustomerRepository) Update(ctx context.Context, id primitive.ObjectID, customer *domain.Customer) error {
	update := bson.M{
		"$set": bson.M{
			"first_name": customer.FirstName,
			"last_name":  customer.LastName,
			"age":        customer.Age,
			"email":      customer.Email,
			"username":   customer.Username,
			"password":   customer.Password,
			"updated_at": customer.UpdatedAt,
		},
	}
	_, err := m.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}
