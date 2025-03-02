package models

type Review struct {
	UserName string `json:"username"`
	Content string `json:"content"`
	Rating int `json:"rating"`
}