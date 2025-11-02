package handlers

import (
	"context"
	"net/http"

	"github.com/SomeSuperCoder/global-chat/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PushRequest[T any] struct {
	Payload []*T `json:"payload" validate:"required"`
	Positon *int `json:"position" validate:"required"`
}

type Pusher[T any] interface {
	Push(ctx context.Context, id bson.ObjectID, values []*T, position int) error
}

func Push[T any, R any](w http.ResponseWriter, r *http.Request, repo Pusher[T], valueGenerator ValueGenerator[T, R]) {
	var request = new(PushRequest[R])

	var id bson.ObjectID
	var exit bool
	if id, exit = utils.ParseRequestID(w, r); exit {
		return
	}

	if DefaultParseAndValidate(w, r, request) {
		return
	}

	transformedPayload := make([]*T, len(request.Payload))
	for i, v := range request.Payload {
		transformedPayload[i] = valueGenerator(v)
	}

	err := repo.Push(r.Context(), id, transformedPayload, *request.Positon)
	if utils.CheckError(w, err, "Failed to push", http.StatusInternalServerError) {
		return
	}
}
