package repository

import "github.com/SomeSuperCoder/global-chat/models"

type SentenceRepo = GenericArrayRepo[models.Sentence, models.SentenceUpdateRequest]

func NewSentenceRepo(parent *PageRepo) *SentenceRepo {
	return ToArrayRepo[models.Sentence, models.SentenceUpdateRequest](parent)
}
