package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Owner     string             `bson:"owner" json:"owner"`
	Balance   int64              `bson:"balance" json:"balance"`
	Currency  string             `bson:"currency" json:"currency"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type Entry struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AccountID primitive.ObjectID `bson:"account_id" json:"account_id"`
	Amount    int64              `bson:"amount" json:"amount"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type Transfer struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	FromAccountID primitive.ObjectID `bson:"fromaccountid"`
	ToAccountID   primitive.ObjectID `bson:"toaccountid"`
	Amount        int64              `bson:"amount"`
	CreatedAt     time.Time          `bson:"createdat"`
	Status        string             `bson:"status"`
}
