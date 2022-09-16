package ports

import "github.com/akpor-kofi/auth/models"

type UserRepository interface {
	Create(user *models.User) error
	Get(id string) (*models.User, error)
	List() ([]*models.User, error)
	Update(id string, user *models.User) (*models.User, error)
	Delete(id string) error
	GetByEmail(email string) (*models.User, error)
}
