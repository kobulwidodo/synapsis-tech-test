package entity

import (
	"go-clean/src/lib/auth"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Password string `json:"-"`
	Name     string
}

type CreateUserParam struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
	Name     string `binding:"required"`
}

type LoginUserParam struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) ConvertToAuthUser() auth.User {
	return auth.User{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
	}
}
