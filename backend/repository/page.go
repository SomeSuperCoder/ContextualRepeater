package repository

import (
	"github.com/SomeSuperCoder/global-chat/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PageRepo = GenericRepo[models.Page, models.PageUpdateRequest]

func NewPageRepo(database *mongo.Database) *PageRepo {
	return NewGenericRepo[models.Page, models.PageUpdateRequest](database, "pages")
}
