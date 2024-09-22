package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"memo/api/notes/models"
)

type FetchFilter struct {
	Count  int
	Sort   string
	UserId string
	Type   string
}

type NotesRepository interface {
	Add(note models.EmbeddedNote, userId string, ctx context.Context) (string, error)
	List(filter FetchFilter, ctx context.Context) ([]*models.BaseNote, error)
	GetById(oId string, userId string, ctx context.Context) (*models.EmbeddedNote, error)
	Update(note *models.EmbeddedNote, ctx context.Context) error
	Delete(id string, userId string, ctx context.Context) error
}

type notesRepository struct {
	client *mongo.Database
}

func NewNotes(client *mongo.Database) NotesRepository {
	return &notesRepository{client}
}

func (r *notesRepository) GetById(oId string, userId string, ctx context.Context) (*models.EmbeddedNote, error) {
	collection := r.client.Collection("notes")

	id, err := primitive.ObjectIDFromHex(oId)
	if err != nil {
		return nil, err
	}

	uId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		primitive.E{Key: "_id", Value: id},
		primitive.E{Key: "user_id", Value: uId},
	}

	var note *models.EmbeddedNote
	err = collection.FindOne(ctx, filter).Decode(&note)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (r *notesRepository) Delete(oId string, userId string, ctx context.Context) error {
	collection := r.client.Collection("notes")

	id, err := primitive.ObjectIDFromHex(oId)
	if err != nil {
		return err
	}

	oUserId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	filter := bson.D{
		primitive.E{Key: "_id", Value: id},
		primitive.E{Key: "user_id", Value: oUserId},
	}

	if _, err = collection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}

func (r *notesRepository) Add(note models.EmbeddedNote, userId string, ctx context.Context) (string, error) {
	collection := r.client.Collection("notes")

	if note.Type == "todo" {
		for i := range note.TodoNote.Tasks {
			note.TodoNote.Tasks[i].ID = primitive.NewObjectID()
		}
	}

	objId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return "", err
	}

	note.UserId = objId
	insertResult, err := collection.InsertOne(ctx, note)

	if err != mongo.ErrNilCursor {
		return "", err
	}

	if oidResult, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
		return oidResult.Hex(), nil
	} else {
		return "", err
	}
}

func (r *notesRepository) List(filter FetchFilter, ctx context.Context) ([]*models.BaseNote, error) {
	sortLayout := 1 // asc
	if filter.Sort == "desc" || filter.Sort == "" {
		sortLayout = -1 // desc
	}

	findOptions := options.Find()
	// findOptions.SetLimit(int64(count))
	findOptions.SetSort(bson.D{{Key: "created_at", Value: sortLayout}})

	collection := r.client.Collection("notes")

	query := bson.D{}
	if filter.UserId != "" {
		objId, err := primitive.ObjectIDFromHex(filter.UserId)
		if err != nil {
			return nil, err
		}

		query = append(query, primitive.E{Key: "user_id", Value: objId})
	}

	if filter.Type != "all" && filter.Type != "" {
		query = append(query, primitive.E{Key: "type", Value: filter.Type})
	}

	cursor, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, err
	}

	note := []*models.BaseNote{}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cursor.Next(ctx) {
		// create a value into which the single document can be decoded
		var elem models.BaseNote
		if err := cursor.Decode(&elem); err != nil {
			return nil, err
		}

		note = append(note, &elem)
	}

	err = cursor.Close(ctx)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (r *notesRepository) Update(note *models.EmbeddedNote, ctx context.Context) error {

	filter := bson.M{"_id": note.ID}

	update := bson.M{"$set": bson.M{}}

	if note.Type == "text" {
		update["$set"].(bson.M)["text_note.content"] = note.TextNote.Content
	} else if note.Type == "movie" {
		update["$set"].(bson.M)["movie_note.year"] = note.MovieNote.Year
		update["$set"].(bson.M)["movie_note.director"] = note.MovieNote.Director
	} else if note.Type == "todo" {
		// TODO: Unimplemented todo update
	}

	update["$set"].(bson.M)["title"] = note.Title
	update["$set"].(bson.M)["updated_at"] = time.Now()

	_, err := r.client.Collection("notes").UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
