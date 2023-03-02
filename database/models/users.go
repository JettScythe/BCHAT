package models

type User struct {
	Username string `binding:"required,min=5,max=30"`
	Password string `binding:"required,min=8"`
}

var Users []*User
