package models

import (
	"memo/api/auth"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Permission string

const (
	PermissionRead  Permission = "read"
	PermissionWrite Permission = "write"
)

type SharedUser struct {
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	Permission Permission         `bson:"permission" json:"permission"`
}

type BaseNote struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type       string             `bson:"type" json:"type"`
	Title      string             `bson:"title" json:"title"`
	Tags       []string           `bson:"tags" json:"tags"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
	UserId     primitive.ObjectID `bson:"user_id,omitempty" json:"user_id"`
	SharedWith []SharedUser       `bson:"shared_with" json:"shared_with,omitempty"`
}

func (n *BaseNote) OwnedBy(user string) bool {
	return user == n.UserId.Hex()
}

// EmbeddedNote uses embedded documents for specific note types
type EmbeddedNote struct {
	BaseNote  `bson:",inline"`
	TextNote  *TextNoteData  `bson:"text_note,omitempty" json:"text_note,omitempty"`
	TodoNote  *TodoNoteData  `bson:"todo_note,omitempty" json:"todo_note,omitempty"`
	MovieNote *MovieNoteData `bson:"movie_note,omitempty" json:"movie_note,omitempty"`
}

type UserNote struct {
	BaseNote   `bson:",inline"`
	SharedWith []struct {
		User       auth.AuthUser `bson:"user" json:"user"`
		Permission Permission    `bson:"permission" json:"permission"`
	} `bson:"shared_with" json:"shared_with,omitempty"`
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
