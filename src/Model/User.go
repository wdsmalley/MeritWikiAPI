package model

import (
)

type User struct {
	UserID int
	UserName string
	Password string
	Email string
	Status int
	Rating int
	Resp UserResponse
}