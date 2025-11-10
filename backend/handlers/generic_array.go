package handlers

import (
	"net/http"

	"github.com/SomeSuperCoder/global-chat/repository"
	"github.com/SomeSuperCoder/global-chat/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ArrayWrapper[T any] struct {
	Payload []T `validate:"dive"`
}

func Push[T, R, IU any](w http.ResponseWriter, r *http.Request, repo *repository.GenericArrayRepo[T, IU], valueGenerator ValueGenerator[T, R], idList []string) {
	var request = ArrayWrapper[*R]{}

	var id bson.ObjectID
	var exit bool
	if id, exit = utils.ParseRequestID(w, r); exit {
		return
	}

	if DefaultParseAndValidate(w, r, &request) {
		return
	}

	transformedPayload := make([]*T, len(request.Payload))
	for i, v := range request.Payload {
		transformedPayload[i] = valueGenerator(v)
	}

	fieldPath, err := repository.FromArrayFieldPath[T, IU](r, idList)
	if utils.CheckError(w, err, "Failed to read array field path", http.StatusBadRequest) {
		return
	}
	err = repo.Push(r.Context(), id, transformedPayload, fieldPath)
	if utils.CheckError(w, err, "Failed to push", http.StatusInternalServerError) {
		return
	}
}

func Pull[IT, IU any](w http.ResponseWriter, r *http.Request, repo *repository.GenericArrayRepo[IT, IU], idList []string) {
	var id bson.ObjectID
	var exit bool
	if id, exit = utils.ParseRequestID(w, r); exit {
		return
	}

	fieldPath, err := repository.FromArrayFieldPath[IT, IU](r, idList)
	if utils.CheckError(w, err, "Failed to read array field path", http.StatusBadRequest) {
		return
	}
	err = repo.Pull(r.Context(), id, fieldPath)
	if utils.CheckError(w, err, "Failed to pull", http.StatusInternalServerError) {
		return
	}
}

func ArrayUpdate[U any, IT any](w http.ResponseWriter, r *http.Request, repo *repository.GenericArrayRepo[IT, U], idList []string) {
	var request = new(U)

	var id bson.ObjectID
	var exit bool
	if id, exit = utils.ParseRequestID(w, r); exit {
		return
	}

	if DefaultParseAndValidate(w, r, request) {
		return
	}

	fieldPath, err := repository.FromArrayFieldPath[IT, U](r, idList)
	if utils.CheckError(w, err, "Failed to read array field path", http.StatusBadRequest) {
		return
	}
	err = repo.ArrayUpdate(r.Context(), id, request, fieldPath)
	if utils.CheckError(w, err, "Failed to update an array element", http.StatusInternalServerError) {
		return
	}
}
