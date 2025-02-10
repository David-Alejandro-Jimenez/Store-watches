package models

type Comment struct {
	ID int `json:"id"`
	Date string `json:"date"`
	UserName string `json:"username"`
	Content string `json:"content"`
	Rating int `json:"rating"`
}