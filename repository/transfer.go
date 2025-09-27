package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"project/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func InitTransferCollection(db *mongo.Database) {
	TransferCollection = db.Collection("transfer")
}

func InitMongoDB(uri string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	// Assign to global variable
	MongoClient = client
}
func InsertOneTransfer(client *mongo.Client, transfer model.Transfer) error {
	session, err := client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.Background())

	txnOpts := options.Transaction().
		SetReadConcern(readconcern.Snapshot()).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 🚫 Prevent self-transfer
		if transfer.FromAccountID == transfer.ToAccountID {
			return nil, errors.New("cannot transfer to the same account")
		}
		if transfer.Amount <= 0 {
			return nil, errors.New("transfer amount must be greater than zero")
		}

		// ✅ Check sender account
		var sender model.Account
		err := AccountCollection.FindOne(sessCtx, bson.M{"_id": transfer.FromAccountID}).Decode(&sender)
		if err != nil {
			return nil, errors.New("sender account not found")
		}
		if sender.Balance < transfer.Amount {
			return nil, errors.New("insufficient balance")
		}

		// ✅ Check receiver account
		var receiver model.Account
		err = AccountCollection.FindOne(sessCtx, bson.M{"_id": transfer.ToAccountID}).Decode(&receiver)
		if err != nil {
			return nil, errors.New("receiver account not found")
		}

		// 🔄 Deduct from sender
		_, err = AccountCollection.UpdateOne(
			sessCtx,
			bson.M{"_id": transfer.FromAccountID},
			bson.M{"$inc": bson.M{"balance": -transfer.Amount}},
		)
		if err != nil {
			return nil, err
		}

		// 🔄 Add to receiver
		_, err = AccountCollection.UpdateOne(
			sessCtx,
			bson.M{"_id": transfer.ToAccountID},
			bson.M{"$inc": bson.M{"balance": transfer.Amount}},
		)
		if err != nil {
			return nil, err
		}

		// 📝 Create debit entry
		senderEntry := model.Entry{
			AccountID: transfer.FromAccountID,
			Amount:    -transfer.Amount,
		}
		if _, err = EntryCollection.InsertOne(sessCtx, senderEntry); err != nil {
			return nil, err
		}

		// 📝 Create credit entry
		receiverEntry := model.Entry{
			AccountID: transfer.ToAccountID,
			Amount:    transfer.Amount,
		}
		if _, err = EntryCollection.InsertOne(sessCtx, receiverEntry); err != nil {
			return nil, err
		}

		// 📌 Add metadata
		transfer.Status = "SUCCESS"
		transfer.CreatedAt = time.Now()

		// 💾 Insert transfer record
		fmt.Printf("Saving transfer: From=%s, To=%s\n", transfer.FromAccountID.Hex(), transfer.ToAccountID.Hex())

		if _, err = TransferCollection.InsertOne(sessCtx, transfer); err != nil {
			return nil, err
		}

		return nil, nil
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)
	return err
}

func UpdateOneTransfer(transfer model.Transfer) (*model.Transfer, error) {
	filter := bson.D{{"id", transfer.ID}}

	update := bson.D{
		{"$set", bson.D{
			{"fromaccountid", transfer.FromAccountID},
			{"toaccountid", transfer.ToAccountID},
			{"amount", transfer.Amount},
		}},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated model.Transfer
	err := TransferCollection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&updated)
	if err != nil {
		log.Fatal(err)
	}
	return &updated, nil

}

func GetAllTransferByID(accountID primitive.ObjectID) ([]model.Transfer, error) {
	var transfers []model.Transfer
	filter := bson.M{
		"$or": []bson.M{
			{"fromaccountid": accountID},
			{"toaccountid": accountID},
		},
	}

	ctx := context.Background()
	cursor, err := TransferCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &transfers); err != nil {
		return nil, err
	}

	// ensure an empty slice instead of nil
	if transfers == nil {
		transfers = []model.Transfer{}
	}

	return transfers, nil
}

type TransferRepository struct {
	Collection *mongo.Collection
}

func GetTransfersWithAccounts(collection *mongo.Collection, ctx context.Context) ([]bson.M, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "accounts"},
			{Key: "localField", Value: "fromaccountid"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "fromAccount"},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "accounts"},
			{Key: "localField", Value: "toaccountid"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "toAccount"},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	// 🔍 Check for missing accounts
	for _, res := range results {
		if len(res["fromAccount"].(bson.A)) == 0 {
			return nil, errors.New("sender account not found")
		}
		if len(res["toAccount"].(bson.A)) == 0 {
			return nil, errors.New("receiver account not found")
		}
	}

	return results, nil
}
