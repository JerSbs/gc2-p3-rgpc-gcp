package infra

import (
	"context"
	"errors"
	"shopping-service/config"
	"shopping-service/internal/shopping/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// interface repository untuk transaksi
type TransactionRepository interface {
	Insert(transaction *domain.Transaction) error
	FindAll() ([]domain.Transaction, error)
	FindByID(id string) (*domain.Transaction, error)
	Update(id string, transaction *domain.Transaction) error
	Delete(id string) error
	DeleteFailedOlderThan(duration time.Duration) error // untuk cron job
}

type transactionRepository struct {
	col *mongo.Collection
}

// inisialisasi collection "transactions"
func NewTransactionRepo() TransactionRepository {
	col := config.DB.Collection("transactions")
	return &transactionRepository{col: col}
}

func (r *transactionRepository) Insert(transaction *domain.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	transaction.CreatedAt = time.Now()
	transaction.Status = "success"

	_, err := r.col.InsertOne(ctx, transaction)
	return err
}

func (r *transactionRepository) FindAll() ([]domain.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []domain.Transaction
	for cursor.Next(ctx) {
		var t domain.Transaction
		if err := cursor.Decode(&t); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (r *transactionRepository) FindByID(id string) (*domain.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var result domain.Transaction
	err = r.col.FindOne(ctx, bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (r *transactionRepository) Update(id string, transaction *domain.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"quantity":   transaction.Quantity,
			"status":     transaction.Status,
			"total":      transaction.Total,
			"product_id": transaction.ProductID,
			"payment_id": transaction.PaymentID,
		},
	}

	_, err = r.col.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func (r *transactionRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// untuk cron job: hapus transaksi failed lebih dari durasi tertentu
func (r *transactionRepository) DeleteFailedOlderThan(duration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cutoff := time.Now().Add(-duration)
	filter := bson.M{
		"status":     "failed",
		"created_at": bson.M{"$lt": cutoff},
	}

	_, err := r.col.DeleteMany(ctx, filter)
	return err
}
