package models

import (
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"unique" json:"email"`
	PasswordHash string `json:"passwordHash"`
}

type UserCreate struct {
	Email    string `validate:"required" json:"email"`
	Password string `validate:"required" json:"password"`
}

type UserSignIn struct {
	Email    string `validate:"required" json:"email"`
	Password string `validate:"required" json:"password"`
}

func (bc *UserCreate) Validate() error {
	validate := validator.New()
	if err := validate.Struct(bc); err != nil {
		return err
	}
	return nil
}

func (us *UserSignIn) Validate() error {
	validate := validator.New()
	if err := validate.Struct(us); err != nil {
		return err
	}
	return nil
}

func (bc *UserCreate) Convert() (*User, error) {
	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(bc.Password), bcrypt.DefaultCost)
	if err != nil {
		return &User{}, err
	}

	return &User{
		Email:        bc.Email,
		PasswordHash: string(hashedPasswd),
	}, nil
}
