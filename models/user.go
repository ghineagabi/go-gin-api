package models

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

type LoginFromHeader struct {
	Auth string `header:"Authorization" binding:"required"`
}

type UserCredentials struct {
	Email string
	Pass  string
}
