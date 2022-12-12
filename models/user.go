package models

type AbstractUser struct {
	FirstName    string `json:"firstName" binding:"required,spacetrim"`
	LastName     string `json:"lastName" binding:"required,spacetrim"`
	isSuperUser  bool
	isNormalUser bool
	isActive     bool
	Password     string `json:"password" binding:"required,pw"`
	Email        string `json:"email" binding:"required,email"`
}

type PartnerUser struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
}

type AbstractUserToUpdate struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type LoginFromHeader struct {
	Auth string `header:"Authorization" binding:"required"`
}

type UserCredentials struct {
	Pass  string `json:"password" binding:"required,pw"`
	Email string `json:"email" binding:"required,email"`
}

type ResetPassword struct {
	Password string `json:"password" binding:"required,pw"`
}
