package models

type Review string

const (
	Pass      Review = "pass"
	Fail      Review = "fail"
	Uncertain Review = "uncertain"
)

type Sentance struct {
	MainContent  string   `json:"main_content"`
	ExtraContent string   `json:"extra_content"`
	Reviews      []Review `json:"reviews"`
}
