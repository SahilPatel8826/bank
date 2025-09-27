package repository

import (
	"context"
	"errors"
	"fmt"
	"project/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Insert user only if account exists
func InsertOneUser(user model.User) error {
	// make sure the collection is initialized
	if AccountCollection == nil || UserCollection == nil {
		return errors.New("collections not initialized")
	}

	// check if account exists directly using repo-level collection
	var account bson.M
	err := AccountCollection.FindOne(context.Background(), bson.M{"_id": user.AccountID}).Decode(&account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("account does not exist, cannot create user")
		}
		return err
	}
	if account["owner"] != user.Name {
		return errors.New("user name does not match account owner name")
	}

	var existingUser model.User
	err = UserCollection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		return errors.New("email already exists")
	} else if err != mongo.ErrNoDocuments {
		return err
	}

	// 4. Check if account already assigned to some user
	err = UserCollection.FindOne(context.Background(), bson.M{"accountid": user.AccountID}).Decode(&existingUser)
	if err == nil {
		return errors.New("this account is already linked to another user")
	} else if err != mongo.ErrNoDocuments {
		return err
	}

	// insert user
	inserted, err := UserCollection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}
	fmt.Println("Inserted user with ID:", inserted.InsertedID)
	return nil
}

func LoginUser(ctx context.Context, email string) (model.User, error) {
	var user model.User
	err := UserCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return user, err
}

func FindAllUser() ([]model.User, error) {
	var fetchedUser []model.User
	cursor, err := UserCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		fetchedUser = append(fetchedUser, user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return fetchedUser, nil
}

func CheckBalance(accountID primitive.ObjectID) (int64, error) {
	var account model.Account
	err := AccountCollection.FindOne(context.TODO(), bson.M{"_id": accountID}).Decode(&account)
	if err != nil {
		return 0, err
	}
	fmt.Println("Account balance:", account.Balance)
	return account.Balance, nil
}

// Fetch user with account check
func GetUserWithAccount(collection *mongo.Collection, ctx context.Context) ([]bson.M, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "accounts"},
			{Key: "localField", Value: "accountid"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "account"},
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

	for _, res := range results {
		accounts := res["account"].(bson.A) // cast directly
		if len(accounts) == 0 {
			return nil, errors.New("invalid accountid: no account found")
		}

		acc := accounts[0].(bson.M) // take first account
		if res["name"] != acc["ownername"] {
			return nil, errors.New("user name and account owner name do not match")
		}
	}

	return results, nil
}
