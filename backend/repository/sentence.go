package repository

import (
	"github.com/SomeSuperCoder/global-chat/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SentenceRepo = GenericArrayRepo[models.Sentence, models.SentenceUpdateRequest]

func NewSentenceRepo(database *mongo.Database) *SentenceRepo {
	return NewGenericArrayRepo[models.Sentence, models.SentenceUpdateRequest](database, "pages", "sentences")
}
