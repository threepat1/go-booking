package repository

import (
	"context"
	"time"

	"booking-app/modules/customer/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (m *MongoCustomerRepository) Save(ctx context.Context, customer *domain.Customer) error {
	_, err := m.collection.InsertOne(ctx, customer)
	return err
}

func (m *MongoCustomerRepository) Update(ctx context.Context, id primitive.ObjectID, customer *domain.Customer) error {
	_, err := m.collection.ReplaceOne(ctx, bson.M{"_id": id}, customer)
	return err
}

func (m *MongoCustomerRepository) VerifyEmail(ctx context.Context, token string) error {
	filter := bson.M{"verification_token": token}
	update := bson.M{
		"$set": bson.M{
			"is_verified": true,
		},
		"$unset": bson.M{
			"verification_token": "",
		},
	}
	result := m.collection.FindOneAndUpdate(ctx, filter, update)
	return result.Err()
}
