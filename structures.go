package main

type AbstractUser struct {
	Age          uint16 `json:"age"`
	FirstName    string `json:"first-name"`
	LastName     string `json:"last-name"`
	isSuperUser  bool
	isNormalUser bool
	isActive     bool
	Password     string `json:"password" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
}

type Post struct {
	Email     string `json:"email"`
	Title     string `json:"title" binding:"spacetrim"`
	GroupName string `json:"group-name"`
	Content   string `json:"content"`
}

type GeneralQueryFields struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

type PostInfo struct {
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
	Content   string `json:"content"`
}

type ToUpdatePost struct {
	Id      int    `json:"id" binding:"required"`
	Title   string `json:"title" binding:"spacetrim"`
	Content string `json:"content"`
}

type AbstractUserSession struct {
	Id string
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
