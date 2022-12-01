package main

import (
	"time"
)

type AbstractUser struct {
	Age          uint16 `json:"age"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	isSuperUser  bool
	isNormalUser bool
	isActive     bool
	Password     string `json:"password" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
}

type AbstractUserToUpdate struct {
	Age       uint16 `json:"age"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type GeneralQueryFields struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

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

type UserLoginFromHeader struct {
	Auth string `header:"Authorization" binding:"required"`
}

type UserCredentials struct {
	Email string
	Pass  string
}

type FileCredentials struct {
	Host               string `json:"host" binding:"required"`
	Port               int    `json:"port" binding:"required"`
	User               string `json:"user" binding:"required"`
	Password           string `json:"password" binding:"required"`
	Dbname             string `json:"dbname" binding:"required"`
	AnonymousGMailName string `json:"anonymousGMailName" binding:"required"`
	AnonymousGmailPass string `json:"anonymousGMailPass" binding:"required"`
}

type VerificationTTL struct {
	AbsUsr AbstractUser
	TTL    time.Time
}

type CachedLoginSessions struct {
	EmailID int
	SessTTL time.Time
}

type GeneralID struct {
	ID int `uri:"ID" binding:"required"`
}

type CommentToCreate struct {
	PostID  int    `json:"postID" binding:"required"`
	Content string `json:"content"`
}

type CommentToGet struct {
	FullName      string    `json:"fullName"`
	PostID        int       `json:"postID"`
	Content       string    `json:"content"`
	Date          time.Time `json:"commentDate"`
	CommentID     int       `json:"commentID"`
	IsEdited      bool      `json:"isEdited"`
	NumberOfLikes int       `json:"NumberOfLikes"`
}
