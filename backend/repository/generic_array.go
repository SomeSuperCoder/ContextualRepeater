package repository

import (
	"context"
	"fmt"

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

type Pusher[T any] interface {
	Push(ctx context.Context, id bson.ObjectID, values []*T, position int) error
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

type Puller interface {
	Pull(ctx context.Context, id bson.ObjectID, position int) error
}

func (r *GenericArrayRepo[T]) Pull(ctx context.Context, id bson.ObjectID, position int) error {
	update := bson.M{
		"$unset": bson.M{
			fmt.Sprintf("%s.%d", r.Field, position): "",
		},
	}

	_, err := r.Collection.UpdateByID(ctx, id, update)
	if err != nil {
		return err
	}

	// Remove the null value created by unset
	update = bson.M{
		"$pull": bson.M{
			r.Field: nil,
		},
	}

	_, err = r.Collection.UpdateByID(ctx, id, update)
	return err
}
