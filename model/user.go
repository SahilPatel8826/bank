package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name" validate:"required"`
	Email     string             `bson:"email" json:"email" validate:"required,email"`
	Password  string             `bson:"password" json:"password" validate:"required,min=6"`
	Role      string             `bson:"role" json:"role" validate:"required,oneof=admin user"`
	AccountID primitive.ObjectID `bson:"accountid,omitempty" json:"accountid" validate:"required"`
}
