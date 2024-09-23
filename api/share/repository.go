package share

import (
	"context"
	"fmt"
	"memo/api/notes/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShareRepository interface {
	List(userId string, ctx context.Context) ([]*models.UserNote, error)
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

func (r *shareRepo) List(userId string, ctx context.Context) ([]*models.UserNote, error) {
	userID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	pipeline := mongo.Pipeline{
		// {{Key: "$match", Value: bson.M{"user_id": userID}}},
		{{Key: "$match", Value: bson.M{
			"shared_with.user_id": userID,
			"owner_id":            bson.M{"$ne": userID},
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "shared_with.user_id",
			"foreignField": "_id",
			"as":           "shared_users",
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":        1,
			"type":       1,
			"tags":       1,
			"title":      1,
			"user_id":    1,
			"created_at": 1,
			"updated_at": 1,
			"shared_with": bson.M{
				"$map": bson.M{
					"input": "$shared_with",
					"as":    "share",
					"in": bson.M{
						"user": bson.M{
							"$arrayElemAt": []interface{}{
								bson.M{"$filter": bson.M{
									"input": "$shared_users",
									"cond":  bson.M{"$eq": []interface{}{"$$this._id", "$$share.user_id"}},
								}},
								0,
							},
						},
						"permission": "$$share.permission",
					},
				},
			},
		}}},
	}

	cursor, err := r.client.Collection("notes").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	notes := []*models.UserNote{}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cursor.Next(ctx) {
		// create a value into which the single document can be decoded
		var elem models.UserNote
		if err := cursor.Decode(&elem); err != nil {
			return nil, err
		}

		notes = append(notes, &elem)
	}

	if err = cursor.Close(ctx); err != nil {
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
