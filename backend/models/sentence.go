package models

type Review string

const (
	Pass      Review = "pass"
	Fail      Review = "fail"
	Uncertain Review = "uncertain"
)

type Sentence struct {
	MainContent  string   `json:"main_content" bson:"main_content"`
	ExtraContent string   `json:"extra_content" bson:"extra_content"`
	Reviews      []Review `json:"reviews" bson:"reviews"`
}

type SentenceCreateRequest struct {
	MainContent  string `json:"main_content" bson:"main_content" validate:"required"`
	ExtraContent string `json:"extra_content" bson:"extra_content"`
}

type SentenceUpdateRequest struct {
	MainContent  string `json:"main_content" bson:"main_content"`
	ExtraContent string `json:"extra_content" bson:"extra_content"`
}
