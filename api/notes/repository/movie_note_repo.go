package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MovieNotesRepository interface {
	Update(noteId string, updates map[string]any, ctx context.Context) error
}

type movieNotesRepository struct {
	client *mongo.Database
}

func NewMovieNotes(client *mongo.Database) MovieNotesRepository {
	return &movieNotesRepository{client}
}

func (r *movieNotesRepository) Update(oId string, updates map[string]any, ctx context.Context) error {
	id, err := primitive.ObjectIDFromHex(oId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}

	update := bson.M{"$set": bson.M{}}

	for key, value := range updates {
		update["$set"].(bson.M)["movie_note."+key] = value
	}

	result, err := r.client.Collection("notes").UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	// check if update is successfull
	if result.MatchedCount == 0 {
		return fmt.Errorf("no documents matched the filter")
	}
	return nil
}
