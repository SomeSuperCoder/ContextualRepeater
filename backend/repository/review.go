package repository

type ReviewRepo = GenericArrayRepo[string, string]

func NewReviewRepo(parent *PageRepo) *ReviewRepo {
	return ToArrayRepo[string, string](parent)
}
