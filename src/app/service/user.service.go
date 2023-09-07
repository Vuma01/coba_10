package service

import (
	"coba_01/src/models"
)

type UserServiceInterface interface {

	//===============>[ CRUD ]<===============\\
	Create(*models.User) error
	GetAll(page int, perPage int) ([]models.User, error)
	FindByID(id string) (*models.User, error)
	FindByEmail(email string) ([]models.User, error)
	Update(user *models.User) error
	Delete(id string) error

	//================>[ AUTH ]<=================\\
	Signup(user *models.User) (string, error)
	Login(userResponse *models.User) (string, error)
}
