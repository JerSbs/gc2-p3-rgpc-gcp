package infra

import (
	"context"
	"errors"
	"payment-service/config"
	"payment-service/internal/payment/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Interface repository
type PaymentRepository interface {
	Insert(ctx context.Context, payment domain.Payment) (*mongo.InsertOneResult, error)
	FindByID(ctx context.Context, id string) (domain.Payment, error)
	DeleteByID(ctx context.Context, id string) (domain.Payment, error)
	FindAll(ctx context.Context) ([]domain.Payment, error)
}

// Implementasi repository
type paymentRepository struct {
	collection *mongo.Collection
}

// Inisialisasi repository
func NewPaymentRepository() PaymentRepository {
	return &paymentRepository{
		collection: config.DB.Collection("payments"),
	}
}

// Simpan data payment ke database
func (r *paymentRepository) Insert(ctx context.Context, payment domain.Payment) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return r.collection.InsertOne(ctx, payment)
}

// Ambil data berdasarkan ID
func (r *paymentRepository) FindByID(ctx context.Context, id string) (domain.Payment, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Payment{}, err
	}

	var result domain.Payment
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Payment{}, nil
		}
		return domain.Payment{}, err
	}

	return result, nil
}

// Hapus data berdasarkan ID
func (r *paymentRepository) DeleteByID(ctx context.Context, id string) (domain.Payment, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Payment{}, err
	}

	var deleted domain.Payment

	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&deleted)
	if err != nil {
		return domain.Payment{}, err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return domain.Payment{}, err
	}

	deleted.Status = "deleted"

	return deleted, nil
}

// Ambil semua data
func (r *paymentRepository) FindAll(ctx context.Context) ([]domain.Payment, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []domain.Payment
	for cursor.Next(ctx) {
		var p domain.Payment
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}
		results = append(results, p)
	}

	return results, nil
}
