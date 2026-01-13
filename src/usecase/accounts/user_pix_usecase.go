package usecase_accounts

import (
	entity_accounts "app/entity/accounts"
	"fmt"

	"github.com/google/uuid"
)

type userPixUseCase struct {
	repo IRepositoryUserPix
}

func NewUserPixUseCase(repo IRepositoryUserPix) IUseCaseUserPix {
	return &userPixUseCase{repo: repo}
}

func (u *userPixUseCase) Create(userPix *entity_accounts.UserPix, ownerID int) error {
	// Enforce ownership at creation time by setting the OwnerID
	// The UserPix entity has `Owner *User`. We might need to handle this.
	// Since we are likely using GORM, if we just set the OwnerID (foreign key), it should be enough if the struct supported it.
	// However, `UserPix` has `Owner *User`.
	// We need to fetch the User or construct a generic User with just ID to satisfy the relationship if GORM is used this way.
	// Or we can just set `Owner: &entity_accounts.User{ID: ownerID}`.

	userPix.OwnerID = ownerID

	if err := u.repo.Create(userPix); err != nil {
		return fmt.Errorf("could not create user pix key")
	}
	return nil
}

func (u *userPixUseCase) GetById(id uuid.UUID, ownerID int) (*entity_accounts.UserPix, error) {
	userPix, err := u.repo.FindByIdAndOwner(id, ownerID)
	if err != nil {
		return nil, fmt.Errorf("pix key not found or access denied")
	}
	return userPix, nil
}

func (u *userPixUseCase) GetAll(ownerID int) ([]*entity_accounts.UserPix, error) {
	return u.repo.GetAllByOwner(ownerID)
}

func (u *userPixUseCase) Delete(id uuid.UUID, ownerID int) error {
	// Verify ownership first
	_, err := u.repo.FindByIdAndOwner(id, ownerID)
	if err != nil {
		return fmt.Errorf("pix key not found or access denied")
	}

	if err := u.repo.Delete(id); err != nil {
		return fmt.Errorf("could not delete pix key")
	}
	return nil
}

func (u *userPixUseCase) Update(userPix *entity_accounts.UserPix, ownerID int) error {
	// Verify ownership first using the ID from the payload (which must be set)
	if userPix.ID == nil {
		return fmt.Errorf("pix id is required")
	}

	existing, err := u.repo.FindByIdAndOwner(*userPix.ID, ownerID)
	if err != nil {
		return fmt.Errorf("pix key not found or access denied")
	}

	// Update allowed fields
	existing.PixKey = userPix.PixKey
	// ... update other fields if any ...

	if err := u.repo.Update(existing); err != nil {
		return fmt.Errorf("could not update pix key")
	}
	return nil
}
