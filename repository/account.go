package repository

import (
	"context"
	"fmt"
	"log"

	"project/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitAccountCollection(db *mongo.Database) {
	AccountCollection = db.Collection("account")
}

func InsertOneAccount(account model.Account) {
	inserted, err := AccountCollection.InsertOne(context.Background(), account)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("inserted", inserted.InsertedID)
}

func UpdateOneAccount(req model.Account) (*model.Account, error) {
	filter := bson.D{{"id", req.ID}}

	update := bson.D{
		{"$set", bson.D{
			{"owner", req.Owner},
			{"balance", req.Balance},
			{"currency", req.Currency},
		}},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedAccount model.Account
	err := AccountCollection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&updatedAccount)
	if err != nil {
		return nil, err
	}

	return &updatedAccount, nil
}

func GetOneAccount(id primitive.ObjectID) (*model.Account, error) {
	filter := bson.D{{"id", id}}

	var fetchedAccount model.Account
	err := AccountCollection.FindOne(context.TODO(), filter).Decode(&fetchedAccount)
	if err != nil {
		return nil, err
	}
	return &fetchedAccount, nil
}

func DeleteOneAccount(id int64) {
	filter := bson.D{{"id", id}}

	result, err := AccountCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("deleted", result.DeletedCount)

}
