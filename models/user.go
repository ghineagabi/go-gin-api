package models

type AbstractUser struct {
	FirstName    string `json:"FirstName" binding:"required,spacetrim"`
	LastName     string `json:"LastName" binding:"required,spacetrim"`
	isSuperUser  bool
	isNormalUser bool
	isActive     bool
	Password     string `json:"Password" binding:"required,pw"`
	Email        string `json:"Email" binding:"required,email"`
}

type PartnerUser struct {
	PhoneNumber string `json:"PhoneNumber" binding:"required"`
}

type AbstractUserToUpdate struct {
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
}

type LoginFromHeader struct {
	Auth string `header:"Authorization" binding:"required"`
}

type UserCredentials struct {
	Password string `json:"Password" binding:"required,pw"`
	Email    string `json:"Email" binding:"required,email"`
}

type ForgotPassword struct {
	Password        string `json:"Password" binding:"required,pw,eqfield=ConfirmPassword"`
	ConfirmPassword string `json:"ConfirmPassword" binding:"required,pw"`
}

type ResetPassword struct {
	OldPassword     string `json:"OldPassword" binding:"required,pw,nefield=Password"`
	Password        string `json:"Password" binding:"required,pw,eqfield=ConfirmPassword"`
	ConfirmPassword string `json:"ConfirmPassword" binding:"required,pw"`
}
