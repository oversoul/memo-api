package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Content     string             `bson:"content" json:"content"`
	IsCompleted bool               `bson:"is_completed" json:"is_completed"`
	CompletedAt *time.Time         `bson:"completed_at" json:"completed_at"`
}
