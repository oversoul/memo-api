package auth

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Token struct {
		Id     primitive.ObjectID `bson:"_id,omitempty"`
		Name   string             `bson:"name"`
		Token  string             `bson:"token"`
		UserId primitive.ObjectID `bson:"user_id"`
	}

	AuthUser struct {
		ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
		Name     string             `json:"name" bson:"name"`
		Email    string             `json:"email" bson:"email"`
		Image    string             `json:"image" bson:"image"`
		Password string             `json:"-" bson:"password"`
	}
)
