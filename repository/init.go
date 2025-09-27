// repository/init.go
package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	MongoClient        *mongo.Client
	AccountCollection  *mongo.Collection
	EntryCollection    *mongo.Collection
	TransferCollection *mongo.Collection
	UserCollection     *mongo.Collection
)

func InitCollections(client *mongo.Client, db *mongo.Database) {
	MongoClient = client
	AccountCollection = db.Collection("account")
	EntryCollection = db.Collection("entry")
	TransferCollection = db.Collection("transfer")
	UserCollection = db.Collection("user")
}
