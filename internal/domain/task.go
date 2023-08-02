package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskStatus int

const (
	TODO = iota
	INPROGRESS
	READY
	COMPLETED
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Description string             `bson:"description,omitempty"`
	Creator     string             `bson:"creator,omitempty"`
	Assignee    string             `bson:"assignee,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty"`
	CompletedAt time.Time          `bson:"completed_at,omitempty"`
	TaskStatus  TaskStatus         `bson:"task_status,omitempty"`
}
