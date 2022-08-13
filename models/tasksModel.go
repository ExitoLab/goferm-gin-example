package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Helpee Interest is the model that inserts a helpee interest objects retrived or inserted into the DB
type Tasks struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       *string            `json:"title" validate:"required,min=2,max=100"`
	Description *string            `json:"description" validate:"required,min=2,max=900000"`
	Task_id     string             `json:"task_id"`
	Created_at  time.Time          `json:"created_at"`
	Updated_at  time.Time          `json:"updated_at"`
}
