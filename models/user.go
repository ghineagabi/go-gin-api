package models

type AbstractUser struct {
	FirstName       string `json:"firstName" binding:"required,spacetrim"`
	LastName        string `json:"lastName" binding:"required,spacetrim"`
	isSuperUser     bool
	isNormalUser    bool
	isActive        bool
	ConfirmPassword string `json:"confirmPassword" binding:"required,pw"`
	Password        string `json:"password" binding:"required,pw,eqfield=ConfirmPassword"`
	Email           string `json:"email" binding:"required,email"`
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
	Password string `json:"password" binding:"required,pw"`
	Email    string `json:"email" binding:"required,email"`
}

type ForgotPassword struct {
	Password        string `json:"password" binding:"required,pw,eqfield=ConfirmPassword"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,pw"`
}

type ResetPassword struct {
	OldPassword     string `json:"oldPassword" binding:"required,pw,nefield=Password"`
	Password        string `json:"password" binding:"required,pw,eqfield=ConfirmPassword"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,pw"`
}
