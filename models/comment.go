package models

import "time"

type CommentToCreate struct {
	PostID  int    `json:"PostID" binding:"required"`
	Content string `json:"Content" binding:"spacetrim"`
}

type CommentToGet struct {
	FullName       string    `json:"FullName"`
	PostID         int       `json:"PostID"`
	Content        string    `json:"Content"`
	Date           time.Time `json:"Date"`
	CommentID      int       `json:"CommentID"`
	IsEdited       bool      `json:"IsEdited"`
	NumberOfLikes  int       `json:"NumberOfLikes"`
	RespondingToID *int      `json:"RespondingToID"`
}
