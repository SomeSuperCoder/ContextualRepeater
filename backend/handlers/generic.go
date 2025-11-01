package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SomeSuperCoder/global-chat/internal/validators"
	"github.com/SomeSuperCoder/global-chat/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type GetteryID[T any] interface {
	GetByID(ctx context.Context, id bson.ObjectID) (T, error)
}

func GetByID[T any](w http.ResponseWriter, r *http.Request, repo GetteryID[T]) {
	var parsedId bson.ObjectID
	var exit bool
	if parsedId, exit = utils.ParseRequestID(w, r); exit {
		return
	}

	value, err := repo.GetByID(r.Context(), parsedId)
	if utils.CheckGetFromDB(w, err) {
		return
	}

	utils.RespondWithJSON(w, value)
}

// ====================
type Finder[T any] interface {
	Find(ctx context.Context) ([]T, error)
}

func Get[T any](w http.ResponseWriter, r *http.Request, repo Finder[T]) {
	cases, err := repo.Find(r.Context())
	if utils.CheckError(w, err, "Failed to get from DB", http.StatusInternalServerError) {
		return
	}

	utils.RespondWithJSON(w, cases)
}

// ====================
type PagedFinder[T any] interface {
	FindPaged(ctx context.Context, page, limit int64) ([]T, int64, error)
}

type PagedResponse[T any] struct {
	Values []T   `json:"values"`
	Count  int64 `json:"count"`
}

func FindPaged[T any](w http.ResponseWriter, r *http.Request, repo PagedFinder[T]) {
	// Get data
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	// Validate
	if page == "" {
		http.Error(w, "No page number provided", http.StatusBadRequest)
		return
	}
	if limit == "" {
		http.Error(w, "No limit number provided", http.StatusBadRequest)
		return
	}

	// Parse
	pageNumber, err := strconv.Atoi(page)
	if utils.CheckError(w, err, "Invalid page number", http.StatusBadRequest) {
		return
	}

	limitNumber, err := strconv.Atoi(limit)
	if utils.CheckError(w, err, "Invalid limit number", http.StatusBadRequest) {
		return
	}

	// Do work
	values, totalCount, err := repo.FindPaged(r.Context(), int64(pageNumber), int64(limitNumber))
	if utils.CheckError(w, err, "Failed to get from DB", http.StatusInternalServerError) {
		return
	}

	// Respond
	utils.RespondWithJSON(w, &PagedResponse[T]{
		Values: values,
		Count:  totalCount,
	})
}

// ====================
type ValueGenerator[T any, R any] = func(r *R) *T
type Callback = func(createdID bson.ObjectID)

type Creatator[T any] interface {
	Create(ctx context.Context, value *T) (bson.ObjectID, error)
}

func Create[T any, R any](w http.ResponseWriter, r *http.Request, repo Creatator[T], valueGenerator ValueGenerator[T, R], callback Callback) {
	var request = new(R)

	if DefaultParseAndValidate(w, r, request) {
		return
	}

	createdID, exit := CreateInner(w, r, repo, valueGenerator(request))
	if exit {
		return
	}
	if callback != nil {
		callback(createdID)
	}
}

func CreateInner[T any](w http.ResponseWriter, r *http.Request, repo Creatator[T], newValue *T) (bson.ObjectID, bool) {
	createdID, err := repo.Create(r.Context(), newValue)
	if utils.CheckError(w, err, "Failed to create", http.StatusInternalServerError) {
		return bson.NilObjectID, true
	}
	return createdID, false
}

// ====================
type ValidatorBuilder = func(r *http.Request, id bson.ObjectID) (validators.Validator, bool)

type Updater[T any] interface {
	Update(ctx context.Context, id bson.ObjectID, update *T) error
}

func UpdateInner[T any](w http.ResponseWriter, r *http.Request, repo Updater[T], id bson.ObjectID, request *T) {
	// Do work
	err := repo.Update(r.Context(), id, request)
	if utils.CheckError(w, err, "Failed to update", http.StatusInternalServerError) {
		return
	}

	// Respond
	fmt.Fprintf(w, "Successfully updated")
}

func DefaultUpdate[R any](w http.ResponseWriter, r *http.Request, repo Updater[R]) {
	Update(w, r, repo, nil)
}

func Update[R any](w http.ResponseWriter, r *http.Request, repo Updater[R], validatorBuilder ValidatorBuilder) {
	var request = new(R)

	var id bson.ObjectID
	var exit bool
	if id, exit = utils.ParseRequestID(w, r); exit {
		return
	}

	if validatorBuilder == nil {
		if DefaultParseAndValidate(w, r, request) {
			return
		}
	} else {
		validator, exit := validatorBuilder(r, id)
		if exit {
			return
		}
		if ParseAndValidate(w, r, validator, request) {
			return
		}
	}

	UpdateInner(w, r, repo, id, request)
}

// ====================
type Deleter interface {
	Delete(ctx context.Context, id bson.ObjectID) error
}

func Delete(w http.ResponseWriter, r *http.Request, repo Deleter) {
	// Load data
	var parsedId bson.ObjectID
	var exit bool
	if parsedId, exit = utils.ParseRequestID(w, r); exit {
		return
	}

	// Do work
	err := repo.Delete(r.Context(), parsedId)
	if utils.CheckError(w, err, "Failed to delete", http.StatusInternalServerError) {
		return
	}

	// Respond
	fmt.Fprintf(w, "Successfully deleted")
}

// ===================================================
// Helpers
// ===================================================
type AccessChecker = func(w http.ResponseWriter, r *http.Request, id bson.ObjectID) bool
type Validator = func(w http.ResponseWriter, r *http.Request)

func ParseAndValidate(w http.ResponseWriter, r *http.Request, validator validators.Validator, request any) bool {
	err := json.NewDecoder(r.Body).Decode(request)
	if utils.CheckJSONError(w, err) {
		return true
	}

	err = validator.ValidateRequest(request)
	if utils.CheckJSONValidError(w, err) {
		return true
	}

	return false
}

func DefaultParseAndValidate(w http.ResponseWriter, r *http.Request, request any) bool {
	return ParseAndValidate(w, r, validators.NewDefaultValidator(), request)
}
