package models

type Page struct {
	Title      string     `json:"title"`
	ReviewDate string     `json:"review_date"`
	Sentences  []Sentance `json:"sentences"`
	Language   string     `json:"language"`
}
