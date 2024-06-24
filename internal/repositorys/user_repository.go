package repositorys

import (
	"errors"
	"fmt"

	"github.com/CesarDelgadoM/generator-reports/internal/models"
	"github.com/CesarDelgadoM/generator-reports/pkg/database"
)

type IUserRepository interface {
	GetEmailById(id uint) (string, error)
}

type UserRepository struct {
	db *database.PostgresDB
}

func NewUserRepository(db *database.PostgresDB) IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) GetEmailById(id uint) (string, error) {
	var user models.User

	err := ur.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return "", errors.New(fmt.Sprint("Get user email failed: ", err))
	}

	return user.Email, nil
}
