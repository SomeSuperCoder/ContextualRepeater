package repository

import (
	"net/http"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func FromArrayFieldPath[T, U any](r *http.Request, idList []string) (*ArrayFieldPath[T, U], error) {
	var result = &ArrayFieldPath[T, U]{
		Pairs: make([]FieldIndexPair, len(idList)),
	}

	for i, id := range idList {
		unparsed := r.PathValue(id)
		parsed, err := strconv.Atoi(unparsed)
		if err != nil {
			return nil, err
		}

		result.Pairs[i] = FieldIndexPair{
			Field: id,
			Index: parsed,
		}
	}

	return result, nil
}

type FieldIndexPair struct {
	Field string
	Index int
}

type ArrayFieldPath[T, U any] struct {
	Pairs []FieldIndexPair
}

func (a *ArrayFieldPath[T, U]) FormPushUpdate(values []*T) bson.M {
	return bson.M{
		"$push": bson.M{
			buildArrayField(a.Pairs, true): bson.M{
				"$each":     values,
				"$position": a.getLastIndex(),
			},
		},
	}
}

func (a *ArrayFieldPath[T, U]) FormUnsetUpdate() bson.M {
	return bson.M{
		"$unset": bson.M{
			buildArrayField(a.Pairs, false): "",
		},
	}
}

func (a *ArrayFieldPath[T, U]) FormPullUpdate() bson.M {
	return bson.M{
		"$pull": bson.M{
			buildArrayField(a.Pairs, true): nil,
		},
	}
}

func (a *ArrayFieldPath[T, U]) FormUpdateUpdate(update *U) bson.M {
	return bson.M{
		"$set": bson.M{
			buildArrayField(a.Pairs, false): update,
		},
	}
}

func (a *ArrayFieldPath[T, U]) getLastIndex() int {
	return a.Pairs[len(a.Pairs)-1].Index
}

func buildArrayField(pairs []FieldIndexPair, skipLastIndex bool) string {
	var builder strings.Builder

	lastIndex := len(pairs) - 1
	for i, pair := range pairs {
		if i != 0 {
			builder.WriteString(".")
		}
		builder.WriteString(pair.Field)
		if skipLastIndex && i == lastIndex {
			break
		}
		builder.WriteString(".")
		builder.WriteString(strconv.Itoa(pair.Index))
	}

	return builder.String()
}
