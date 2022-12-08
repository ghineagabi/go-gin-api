package models

type PostToCreate struct {
	Title   string `json:"title" binding:"required,spacetrim"`
	Content string `json:"content" binding:"required,spacetrim"`
}

type PostToGet struct {
	Title    string `json:"title"`
	FullName string `json:"fullName"`
	Content  string `json:"content"`
}

type PostToUpdate struct {
	Id      int    `json:"id" binding:"required"`
	Title   string `json:"title" binding:"spacetrim"`
	Content string `json:"content"`
}
