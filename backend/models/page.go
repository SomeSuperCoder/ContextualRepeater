package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Page struct {
	ID          bson.ObjectID `json:"_id" bson:"_id,omitempty"`
	Title       string        `json:"title" bson:"title"`
	ReviewDate  time.Time     `json:"review_date" bson:"review_date"`
	ReviewsDone uint          `json:"reviews_done" bson:"reviews_done"`
	Sentences   []Sentence    `json:"sentences" bson:"sentences"`
	Language    string        `json:"language" bson:"language"`
}

type PageCreateRequest struct {
	Title    string `json:"title" validate:"required"`
	Language string `json:"language" validate:"required"`
}

type PageUpdateRequest struct {
	Title       *string    `json:"title,omitempty" bson:"title"`
	ReviewDate  *time.Time `json:"review_date,omitempty" bson:"review_date"`
	ReviewsDone *uint      `json:"reviews_done,omitempty" bson:"reviews_done"`
}
