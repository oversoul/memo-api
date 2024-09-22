package share

import (
	"context"
	"fmt"
	"memo/api/notes/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ShareRepository interface {
	List(userId string, ctx context.Context) ([]*models.BaseNote, error)
	ShareNote(request *shareRequest, ctx context.Context) error
}

type shareRepo struct {
	client *mongo.Database
}

type shareRequest struct {
	UserID     string            `json:"user_id" validate:"required"`
	NoteID     string            `json:"note_id" validate:"required"`
	Permission models.Permission `json:"permission" validate:"required|in:read,write"`
}

func NewShareRepo(client *mongo.Database) ShareRepository {
	return &shareRepo{client}
}

func (r *shareRepo) List(userId string, ctx context.Context) ([]*models.BaseNote, error) {
	userID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"$or": []bson.M{
			{"shared_with.user_id": userID},
		},
	}

	options := options.Find().SetProjection(bson.M{
		"type":       1,
		"tags":       1,
		"title":      1,
		"user_id":    1,
		"updated_at": 1,
		"created_at": 1,
		"shared_with": bson.M{
			"$filter": bson.M{
				"input": "$shared_with",
				"as":    "share",
				"cond":  bson.M{"$eq": []interface{}{"$$share.user_id", userID}},
			},
		},
	})

	cursor, err := r.client.Collection("notes").Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}

	notes := []*models.BaseNote{}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cursor.Next(ctx) {
		// create a value into which the single document can be decoded
		var elem models.BaseNote
		if err := cursor.Decode(&elem); err != nil {
			return nil, err
		}

		notes = append(notes, &elem)
	}

	err = cursor.Close(ctx)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *shareRepo) ShareNote(request *shareRequest, ctx context.Context) error {
	noteID, err := primitive.ObjectIDFromHex(request.NoteID)
	if err != nil {
		return err
	}

	userID, err := primitive.ObjectIDFromHex(request.UserID)
	if err != nil {
		return err
	}

	var note *models.BaseNote
	filter := bson.M{"_id": noteID}

	if err = r.client.Collection("notes").FindOne(ctx, filter).Decode(&note); err != nil {
		return err
	}

	// Check if the user already has access
	for _, sharedUser := range note.SharedWith {
		if sharedUser.UserID == userID {
			return fmt.Errorf("User already has access to this note")
		}
	}

	update := bson.M{
		"$addToSet": bson.M{
			"shared_with": models.SharedUser{
				UserID:     userID,
				Permission: request.Permission,
			},
		},
	}

	_, err = r.client.Collection("notes").UpdateOne(context.Background(), filter, update)

	return err
}
