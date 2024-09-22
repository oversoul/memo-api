package repository

import (
	"context"
	"fmt"
	"memo/api/notes/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TodoNotesRepository interface {
	Create(noteId, content string, ctx context.Context) (string, error)
	Update(noteId, taskId string, updates map[string]any, ctx context.Context) error
}

type todoNotesRepository struct {
	client *mongo.Database
}

func NewTodoNotes(client *mongo.Database) TodoNotesRepository {
	return &todoNotesRepository{client}
}

func (r *todoNotesRepository) Update(oId string, taskId string, updates map[string]any, ctx context.Context) error {
	id, err := primitive.ObjectIDFromHex(oId)
	if err != nil {
		return err
	}

	taskPr, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id, "todo_note.tasks._id": taskPr}

	update := bson.M{"$set": bson.M{}}

	for key, value := range updates {
		update["$set"].(bson.M)["todo_note.tasks.$."+key] = value
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

func (r *todoNotesRepository) Create(todoId, content string, ctx context.Context) (string, error) {
	task := models.Task{
		ID:          primitive.NewObjectID(),
		Content:     content,
		IsCompleted: false,
		CompletedAt: nil,
	}

	id, err := primitive.ObjectIDFromHex(todoId)
	if err != nil {
		return "", err
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$push": bson.M{"todo_note.tasks": task}}
	opts := options.Update().SetUpsert(true)

	result, err := r.client.Collection("notes").UpdateOne(ctx, filter, update, opts)

	if result.ModifiedCount > 0 {
		return task.ID.Hex(), nil
	} else {
		return "", fmt.Errorf("The task was not inserted")
	}
}
