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

// struct repo
type ProductRepo struct {
	collection *mongo.Collection
}

// init repo
func NewProductRepo() *ProductRepo {
	collection := config.DB.Collection("products")
	return &ProductRepo{collection: collection}
}

// simpan product
func (r *ProductRepo) CreateProduct(product *domain.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	product.CreatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		product.ID = oid
	} else {
		return errors.New("failed to convert inserted ID to ObjectID")
	}

	return nil
}

// ambil semua product
func (r *ProductRepo) GetAllProducts() ([]domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []domain.Product
	for cursor.Next(ctx) {
		var p domain.Product
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if len(products) == 0 {
		return nil, errors.New("no products found")
	}

	return products, nil
}

// ambil product by ID
func (r *ProductRepo) GetProductByID(id string) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}

	var product domain.Product
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &product, nil
}

// update product by ID
func (r *ProductRepo) UpdateProduct(id string, product *domain.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid product ID")
	}

	update := bson.M{
		"$set": bson.M{
			"name":  product.Name,
			"price": product.Price,
			"stock": product.Stock,
		},
	}

	result, err := r.collection.UpdateByID(ctx, objID, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("product not found")
	}

	return nil
}

// delete product by ID
func (r *ProductRepo) DeleteProduct(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid product ID")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("product not found")
	}

	return nil
}
