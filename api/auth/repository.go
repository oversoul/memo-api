package auth

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepository interface {
	InsertToken(id, name, token string, ctx context.Context) (string, error)
	FindUserById(id string, ctx context.Context) (*AuthUser, error)
	FindUserByEmail(email string, ctx context.Context) (*AuthUser, error)
	FindToken(key string, ctx context.Context) (*Token, error)
	DeleteToken(id string, ctx context.Context) error
	InsertUser(user AuthUser, ctx context.Context) error
	DeleteUserTokens(userId string, ctx context.Context) error

	UpdateUserInfo(u *AuthUser, ctx context.Context) error
}

type authRepository struct {
	client *mongo.Database
}

func NewRepo(client *mongo.Database) AuthRepository {
	return &authRepository{client}
}

func (r *authRepository) InsertUser(user AuthUser, ctx context.Context) error {
	collection := r.client.Collection("users")

	insertResult, err := collection.InsertOne(ctx, user)
	if err != mongo.ErrNilCursor {
		return err
	}

	if _, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
		return nil
	} else {
		return err
	}
}

func (r *authRepository) UpdateUserInfo(user *AuthUser, ctx context.Context) error {
	filter := bson.M{"_id": user.ID}

	update := bson.M{"$set": bson.M{}}

	update["$set"].(bson.M)["name"] = user.Name
	update["$set"].(bson.M)["image"] = user.Image

	_, err := r.client.Collection("users").UpdateOne(ctx, filter, update)
	return err
}

func (r *authRepository) InsertToken(id, name, token string, ctx context.Context) (string, error) {
	collection := r.client.Collection("access_tokens")

	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	model := Token{Name: name, Token: token, UserId: oId}

	insertResult, err := collection.InsertOne(ctx, model)
	if err != nil {
		return "", err
	}

	if oidResult, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
		return oidResult.Hex(), nil
	} else {
		return "", err
	}
}

func (r *authRepository) FindUserById(id string, ctx context.Context) (*AuthUser, error) {
	collection := r.client.Collection("users")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: oid}}

	var user *AuthUser
	if err := collection.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *authRepository) FindUserByEmail(email string, ctx context.Context) (*AuthUser, error) {
	collection := r.client.Collection("users")

	filter := bson.D{primitive.E{Key: "email", Value: email}}

	var user *AuthUser
	if err := collection.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *authRepository) FindToken(key string, ctx context.Context) (*Token, error) {
	collection := r.client.Collection("access_tokens")

	id, err := primitive.ObjectIDFromHex(key)
	if err != nil {
		return nil, err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	var token *Token
	if err := collection.FindOne(ctx, filter).Decode(&token); err != nil {
		return nil, err
	}

	return token, nil
}

func (r *authRepository) DeleteToken(id string, ctx context.Context) error {
	collection := r.client.Collection("access_tokens")

	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: oId}}
	if _, err = collection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}

func (r *authRepository) DeleteUserTokens(userId string, ctx context.Context) error {
	collection := r.client.Collection("access_tokens")

	oId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	filter := bson.D{primitive.E{Key: "user_id", Value: oId}}
	if _, err = collection.DeleteMany(ctx, filter); err != nil {
		return err
	}

	return nil
}
