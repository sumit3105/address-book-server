package repository

import (
	"address-book-server/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (repo *userRepository) Create(user *model.User) error {
	return repo.db.Create(user).Error
}

func (repo *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User

	err := repo.db.Where("email = ? AND is_deleted = false", email).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}