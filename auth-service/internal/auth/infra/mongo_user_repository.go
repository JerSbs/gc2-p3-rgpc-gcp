package infra

import (
	"context"
	"errors"
	"time"

	"auth-service/internal/auth/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(db *mongo.Database) *MongoUserRepository {
	return &MongoUserRepository{
		collection: db.Collection("users"),
	}
}

func (r *MongoUserRepository) Save(user domain.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(context.Background(), user)
	return err
}

func (r *MongoUserRepository) GetByEmail(email string) (domain.User, error) {
	var user domain.User
	filter := bson.M{"email": email}
	err := r.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}
	return user, nil
}
