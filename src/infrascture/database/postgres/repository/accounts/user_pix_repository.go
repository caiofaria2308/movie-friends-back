package repository_accounts

import (
	entity_accounts "app/entity/accounts"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userPixRepository struct {
	DB *gorm.DB
}

func NewUserPixRepository(db *gorm.DB) *userPixRepository {
	return &userPixRepository{DB: db}
}

func (r *userPixRepository) Create(userPix *entity_accounts.UserPix) error {
	return r.DB.Create(userPix).Error
}

func (r *userPixRepository) FindByIdAndOwner(idData uuid.UUID, ownerID int) (*entity_accounts.UserPix, error) {
	var userPix entity_accounts.UserPix
	// We need to join with User or check the OwnerID foreign key if it existed directly,
	// but based on struct, UserPix has `Owner *User`. GORM usually creates `owner_id`.
	// Let's assume `owner_id` is the column name for the `Owner` relationship.
	// Since `User` struct has `ID int`, the foreign key is likely `owner_id` (int).

	if err := r.DB.Preload("Owner").Where("id = ? AND owner_id = ?", idData, ownerID).First(&userPix).Error; err != nil {
		return nil, err
	}
	return &userPix, nil
}

func (r *userPixRepository) GetAllByOwner(ownerID int) ([]*entity_accounts.UserPix, error) {
	var userPixs []*entity_accounts.UserPix
	if err := r.DB.Where("owner_id = ?", ownerID).Find(&userPixs).Error; err != nil {
		return nil, err
	}
	return userPixs, nil
}

func (r *userPixRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&entity_accounts.UserPix{}, "id = ?", id).Error
}

func (r *userPixRepository) Update(userPix *entity_accounts.UserPix) error {
	return r.DB.Save(userPix).Error
}
