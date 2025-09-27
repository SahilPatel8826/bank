package repository

import (
	"context"
	"fmt"
	"log"
	"project/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitEntryCollection(db *mongo.Database) {
	EntryCollection = db.Collection("entry")
}

func CreateEntry(account model.Entry) {
	inserted, err := EntryCollection.InsertOne(context.Background(), account)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("inserted enquiry", inserted.InsertedID)

}

func UpdateOneEntry(req model.Entry) (*model.Entry, error) {
	filter := bson.D{{"_id", req.ID}}

	update := bson.D{
		{"$set", bson.D{
			{"accountid", req.AccountID},
			{"amount", req.Amount},
		}},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedEntry model.Entry
	err := EntryCollection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&updatedEntry)
	if err != nil {
		return nil, err
	}

	return &updatedEntry, nil
}
