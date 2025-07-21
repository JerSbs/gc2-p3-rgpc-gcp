package domain

import "time"

// struct untuk model Transaction
type Transaction struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	ProductID string    `bson:"product_id" json:"product_id"`
	PaymentID string    `bson:"payment_id" json:"payment_id"`
	Email     string    `bson:"email" json:"email"`
	Quantity  int       `bson:"quantity" json:"quantity"`
	Total     float64   `bson:"total" json:"total"`
	Status    string    `bson:"status" json:"status"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
