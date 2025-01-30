package models

type Comment struct {
	Date string `json:"date"`
	Name string `json:"name"`
	Content string `json:"content"`
	Stars int `json:"stars"`
}