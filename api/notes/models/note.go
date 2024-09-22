package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseNote struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type      string             `bson:"type" json:"type"`
	Title     string             `bson:"title" json:"title"`
	Tags      []string           `bson:"tags" json:"tags"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	UserId    primitive.ObjectID `bson:"user_id,omitempty" json:"user_id"`
}

// EmbeddedNote uses embedded documents for specific note types
type EmbeddedNote struct {
	BaseNote  `bson:",inline"`
	TextNote  *TextNoteData  `bson:"text_note,omitempty" json:"text_note,omitempty"`
	TodoNote  *TodoNoteData  `bson:"todo_note,omitempty" json:"todo_note,omitempty"`
	MovieNote *MovieNoteData `bson:"movie_note,omitempty" json:"movie_note,omitempty"`
}

type TextNoteData struct {
	Content string `bson:"content" json:"content"`
}

type TodoNoteData struct {
	Tasks []Task `bson:"tasks" json:"tasks"`
}

type MovieNoteData struct {
	Year     int    `bson:"year" json:"year"`
	Watched  bool   `bson:"watched" json:"watched"`
	Director string `bson:"director" json:"director"`
}
