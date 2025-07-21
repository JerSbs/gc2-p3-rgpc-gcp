package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email" json:"email"`
	Amount    float64            `bson:"amount" json:"amount"`
	Status    string             `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
