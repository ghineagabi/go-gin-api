package utils

import (
	"example/web-service-gin/models"
	"time"
)

type GeneralQueryFields struct {
	Limit  int `form:"Limit"`
	Offset int `form:"Offset"`
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

type CachedLoginSessions struct {
	EmailID int
	SessTTL time.Time
}

type GeneralID struct {
	ID int `uri:"ID" binding:"required"`
}

type VerificationTTL struct {
	AbsUsr models.AbstractUser
	TTL    time.Time
}

type UserLoginFromHeader struct {
	Auth string `header:"Authorization" binding:"required"`
}

type GeneralEmail struct {
	Email string `json:"Email" binding:"required,email"`
}
