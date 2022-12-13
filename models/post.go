package models

type PostToCreate struct {
	Title   string `json:"Title" binding:"required,spacetrim"`
	Content string `json:"Content" binding:"required,spacetrim"`
}

type PostToGet struct {
	Title    string `json:"Title"`
	FullName string `json:"FullName"`
	Content  string `json:"Content"`
}

type PostToUpdate struct {
	Id      int    `json:"Id" binding:"required"`
	Title   string `json:"Title" binding:"spacetrim"`
	Content string `json:"Content"`
}
