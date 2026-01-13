package repository_accounts

import (
	entity_accounts "app/entity/accounts"

	"gorm.io/gorm"
)

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) Create(user *entity_accounts.User) error {
	return r.DB.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*entity_accounts.User, error) {
	var user entity_accounts.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindById(id int) (*entity_accounts.User, error) {
	var user entity_accounts.User
	if err := r.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *entity_accounts.User) error {
	return r.DB.Save(user).Error
}
