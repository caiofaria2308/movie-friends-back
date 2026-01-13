package usecase_accounts

import (
	entity_accounts "app/entity/accounts"
	"time"

	"github.com/google/uuid"
)

type IRepositoryUserDayOff interface {
	Create(dayOff *entity_accounts.UserDayOff) error
	CreateBatch(dayOffs []*entity_accounts.UserDayOff) error
	FindByIdAndOwner(id uuid.UUID, ownerID int) (*entity_accounts.UserDayOff, error)
	FindAllByOwner(ownerID int) ([]*entity_accounts.UserDayOff, error)
	FindAllByOwnerWithFilter(ownerID int, startDate, endDate *time.Time) ([]*entity_accounts.UserDayOff, error)
	FindFutureByName(fatherID uuid.UUID, fromDate time.Time, ownerID int) ([]*entity_accounts.UserDayOff, error)
	FindAllByFather(fatherID uuid.UUID, ownerID int) ([]*entity_accounts.UserDayOff, error)
	DeleteById(id uuid.UUID) error
	DeleteBatch(ids []uuid.UUID) error
	Update(dayOff *entity_accounts.UserDayOff) error
}

type IUseCaseUserDayOff interface {
	Create(dayOff *entity_accounts.UserDayOff, ownerID int) error
	Update(dayOff *entity_accounts.UserDayOff, ownerID int, mode string) error
	Delete(id uuid.UUID, ownerID int, mode string) error
	GetById(id uuid.UUID, ownerID int) (*entity_accounts.UserDayOff, error)
	GetAll(ownerID int, filterType string, year, week, month int) ([]*entity_accounts.UserDayOff, error)
}
