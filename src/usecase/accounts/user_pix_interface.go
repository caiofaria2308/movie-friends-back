package usecase_accounts

import (
	entity_accounts "app/entity/accounts"

	"github.com/google/uuid"
)

type IRepositoryUserPix interface {
	Create(userPix *entity_accounts.UserPix) error
	FindByIdAndOwner(id uuid.UUID, ownerID int) (*entity_accounts.UserPix, error)
	GetAllByOwner(ownerID int) ([]*entity_accounts.UserPix, error)
	Delete(id uuid.UUID) error
	Update(userPix *entity_accounts.UserPix) error
}

type IUseCaseUserPix interface {
	Create(userPix *entity_accounts.UserPix, ownerID int) error
	GetById(id uuid.UUID, ownerID int) (*entity_accounts.UserPix, error)
	GetAll(ownerID int) ([]*entity_accounts.UserPix, error)
	Delete(id uuid.UUID, ownerID int) error
	Update(userPix *entity_accounts.UserPix, ownerID int) error
}
