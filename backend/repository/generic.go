package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type GenericRepo[T any, U any] struct {
	database   *mongo.Database
	Collection *mongo.Collection
}

func NewGenericRepo[T any, U any](database *mongo.Database, collectionName string) *GenericRepo[T, U] {
	return &GenericRepo[T, U]{
		database:   database,
		Collection: database.Collection(collectionName),
	}
}

func ToArrayRepo[T any, U any](r *GenericRepo[any, any], fieldPath ArrayFieldPath) *GenericArrayRepo[T, U] {
	return &GenericArrayRepo[T, U]{
		GenericRepo: r,
		FieldPath:   fieldPath,
	}
}

type PagedFinder[T any] interface {
	FindPaged(ctx context.Context, page, limit int64) ([]T, int64, error)
}

func (r *GenericRepo[T, U]) FindPaged(ctx context.Context, page, limit int64) ([]T, int64, error) {
	var values = []T{}

	// Set pagination options
	skip := (page - 1) * limit
	opts := options.Find()
	opts.SetLimit(limit)
	opts.SetSkip(skip)
	opts.SetSort(bson.M{"created_at": -1})

	// Init a cursor
	cursor, err := r.Collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Extract records
	err = cursor.All(ctx, &values)
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	count, err := r.Collection.CountDocuments(ctx, bson.M{})

	return values, count, err
}

type Finder[T any] interface {
	Find(ctx context.Context) ([]T, error)
}

func (r *GenericRepo[T, U]) Find(ctx context.Context) ([]T, error) {
	return FindWithFilter[T](ctx, r.Collection, struct{}{})
}

type Creatator[T any] interface {
	Create(ctx context.Context, value *T) (bson.ObjectID, error)
}

func (r *GenericRepo[T, U]) Create(ctx context.Context, value *T) (bson.ObjectID, error) {
	res, err := r.Collection.InsertOne(ctx, value)
	objID, _ := res.InsertedID.(bson.ObjectID)
	return objID, err
}

type GetterByID[T any] interface {
	GetByID(ctx context.Context, id bson.ObjectID) (T, error)
}

func (r *GenericRepo[T, U]) GetByID(ctx context.Context, id bson.ObjectID) (*T, error) {
	return GetBy[T](ctx, r.Collection, "_id", id)

}

type Updater[T any] interface {
	Update(ctx context.Context, id bson.ObjectID, update *T) error
}

func (r *GenericRepo[T, U]) Update(ctx context.Context, id bson.ObjectID, update *U) error {
	res := r.Collection.FindOneAndUpdate(ctx, bson.M{
		"_id": id,
	}, bson.M{
		"$set": update,
	})

	return res.Err()
}

type Deleter interface {
	Delete(ctx context.Context, id bson.ObjectID) error
}

func (r *GenericRepo[T, U]) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{
		"_id": id,
	})
	if err != nil {
		return err
	}

	return nil
}

// ================================
// ===== Array related stuff ======
// ================================

type GenericArrayRepo[T any, U any] struct {
	*GenericRepo[any, any]
	FieldPath ArrayFieldPath
}

type Pusher[T any] interface {
	Push(ctx context.Context, id bson.ObjectID, values []*T, position int) error
}

func (r *GenericArrayRepo[T, U]) Push(ctx context.Context, id bson.ObjectID, values []*T, position int) error {
	update := bson.M{
		"$push": bson.M{
			r.FieldPath.GetPushPath(): bson.M{
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

func (r *GenericArrayRepo[T, U]) Pull(ctx context.Context, id bson.ObjectID) error {
	update := bson.M{
		"$unset": bson.M{
			r.FieldPath.GetPullPath(): "",
		},
	}

	_, err := r.Collection.UpdateByID(ctx, id, update)
	if err != nil {
		return err
	}

	// Remove the null value created by unset
	update = bson.M{
		"$pull": bson.M{
			r.FieldPath.GetUpdatePath(): nil,
		},
	}

	_, err = r.Collection.UpdateByID(ctx, id, update)
	return err
}

type ArrayUpdater[U any] interface {
	UpdateByIndex(ctx context.Context, id bson.ObjectID, position int, update U) error
}

func (r *GenericArrayRepo[T, U]) Update(ctx context.Context, id bson.ObjectID, update U) error {
	mongoUpdate := bson.M{
		"$set": bson.M{
			r.FieldPath.GetUpdatePath(): update,
		},
	}

	_, err := r.Collection.UpdateByID(ctx, id, mongoUpdate)
	return err
}

// =============================
// ======= Helper funcs ========
// =============================
func FindWithFilter[T any](ctx context.Context, c *mongo.Collection, filter any) ([]T, error) {
	var values = []T{}

	// Init a cursor
	cursor, err := c.Find(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Extract records
	err = cursor.All(ctx, &values)
	if err != nil {
		return nil, err
	}

	return values, err
}

func GetBy[T any](ctx context.Context, c *mongo.Collection, key string, value any) (*T, error) {
	var got T

	err := c.FindOne(ctx, bson.M{
		key: value,
	}, nil).Decode(&got)
	if err != nil {
		return nil, err
	}

	return &got, err
}
