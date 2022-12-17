package models

import "time"

type CommentToCreate struct {
	PostID  int    `json:"postId" binding:"required"`
	Content string `json:"content" binding:"spacetrim"`
}

type CommentToGet struct {
	FullName       string    `json:"fullName"`
	PostID         int       `json:"postId"`
	Content        string    `json:"content"`
	Date           time.Time `json:"date"`
	CommentID      int       `json:"commentID"`
	IsEdited       bool      `json:"isEdited"`
	NumberOfLikes  int       `json:"numberOfLikes"`
	RespondingToID *int      `json:"respondingToID"`
}
