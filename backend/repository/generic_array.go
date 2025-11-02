package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type GenericArrayRepo[T any] struct {
	database   *mongo.Database
	Collection *mongo.Collection
	Field      string
}

func NewGenericArrayRepo[T any](database *mongo.Database, collectionName string, field string) *GenericArrayRepo[T] {
	return &GenericArrayRepo[T]{
		database:   database,
		Collection: database.Collection(collectionName),
		Field:      field,
	}
}

func (r *GenericArrayRepo[T]) Push(ctx context.Context, id bson.ObjectID, values []*T, position int) error {
	update := bson.M{
		"$push": bson.M{
			r.Field: bson.M{
				"$each":     values,
				"$position": position,
			},
		},
	}

	_, err := r.Collection.UpdateByID(ctx, id, update)
	return err
}
