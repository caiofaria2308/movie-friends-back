package repository_accounts

import (
	entity_accounts "app/entity/accounts"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userDayOffRepository struct {
	DB *gorm.DB
}

func NewUserDayOffRepository(db *gorm.DB) *userDayOffRepository {
	return &userDayOffRepository{DB: db}
}

func (r *userDayOffRepository) Create(dayOff *entity_accounts.UserDayOff) error {
	return r.DB.Create(dayOff).Error
}

func (r *userDayOffRepository) CreateBatch(dayOffs []*entity_accounts.UserDayOff) error {
	return r.DB.CreateInBatches(dayOffs, 100).Error
}

func (r *userDayOffRepository) FindByIdAndOwner(id uuid.UUID, ownerID int) (*entity_accounts.UserDayOff, error) {
	var dayOff entity_accounts.UserDayOff
	if err := r.DB.Where("id = ? AND owner_id = ?", id, ownerID).First(&dayOff).Error; err != nil {
		return nil, err
	}
	return &dayOff, nil
}

func (r *userDayOffRepository) FindAllByOwner(ownerID int) ([]*entity_accounts.UserDayOff, error) {
	var dayOffs []*entity_accounts.UserDayOff
	if err := r.DB.Where("owner_id = ?", ownerID).Find(&dayOffs).Error; err != nil {
		return nil, err
	}
	return dayOffs, nil
}

func (r *userDayOffRepository) FindAllByOwnerWithFilter(ownerID int, startDate, endDate *time.Time) ([]*entity_accounts.UserDayOff, error) {
	var dayOffs []*entity_accounts.UserDayOff
	query := r.DB.Where("owner_id = ?", ownerID)

	if startDate != nil && endDate != nil {
		// Filter day-offs that overlap with the date range
		// A day-off overlaps if: init_hour < endDate AND end_hour > startDate
		query = query.Where("init_hour < ? AND end_hour > ?", endDate, startDate)
	}

	if err := query.Order("init_hour ASC").Find(&dayOffs).Error; err != nil {
		return nil, err
	}
	return dayOffs, nil
}

func (r *userDayOffRepository) FindFutureByName(fatherID uuid.UUID, fromDate time.Time, ownerID int) ([]*entity_accounts.UserDayOff, error) {
	var dayOffs []*entity_accounts.UserDayOff
	// Check init_hour >= fromDate. Using InitHour which is *time.Time.
	if err := r.DB.Where("owner_id = ? AND day_off_father_id = ? AND init_hour >= ?", ownerID, fatherID, fromDate).Find(&dayOffs).Error; err != nil {
		return nil, err
	}
	return dayOffs, nil
}

func (r *userDayOffRepository) FindAllByFather(fatherID uuid.UUID, ownerID int) ([]*entity_accounts.UserDayOff, error) {
	var dayOffs []*entity_accounts.UserDayOff
	if err := r.DB.Where("owner_id = ? AND day_off_father_id = ?", ownerID, fatherID).Find(&dayOffs).Error; err != nil {
		return nil, err
	}
	return dayOffs, nil
}

func (r *userDayOffRepository) DeleteById(id uuid.UUID) error {
	return r.DB.Delete(&entity_accounts.UserDayOff{}, "id = ?", id).Error
}

func (r *userDayOffRepository) DeleteBatch(ids []uuid.UUID) error {
	return r.DB.Delete(&entity_accounts.UserDayOff{}, "id IN ?", ids).Error
}

func (r *userDayOffRepository) Update(dayOff *entity_accounts.UserDayOff) error {
	return r.DB.Save(dayOff).Error
}
