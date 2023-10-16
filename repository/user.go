package repository

import (
	"errors"
	"strings"

	"github.com/akhil-is-watching/authservice/models"
	"gorm.io/gorm"
)

type UserRepositry struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepositry {
	return &UserRepositry{db}
}

func (repo *UserRepositry) Create(u *models.User) error {
	if err := repo.db.Create(u).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique") {
			return errors.New("email already exists")
		}
		return err
	}
	return nil
}

func (repo *UserRepositry) Get(email string) (models.User, error) {
	var user models.User
	if err := repo.db.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (repo *UserRepositry) All() ([]models.User, error) {
	var users []models.User
	if err := repo.db.Find(&users).Error; err != nil {
		return users, err
	}
	return users, nil
}
