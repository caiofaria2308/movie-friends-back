package entity_accounts

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	RepeatTypeDaily   = "daily"
	RepeatTypeWeekly  = "weekly"
	RepeatTypeMonthly = "monthly"
	RepeatTypeYearly  = "yearly"
)

type UserDayOff struct {
	ID             *uuid.UUID  `json:"id"`
	InitHour       *time.Time  `json:"init_hour"`
	EndHour        *time.Time  `json:"end_hour"`
	Owner          *User       `json:"owner"`
	OwnerID        int         `json:"owner_id"` // Explicit FK for easier queries
	Repeat         bool        `json:"repeat"`
	RepeatType     string      `json:"repeat_type"`
	RepeatValue    string      `json:"repeat_value"`
	DayOffFatherID *uuid.UUID  `json:"day_off_father_id"`
	DayOffFather   *UserDayOff `json:"day_off_father" gorm:"foreignKey:DayOffFatherID"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

func (c *UserDayOff) TableName() string {
	return "accounts_user_day_off"
}

func (c *UserDayOff) BeforeCreate(tx *gorm.DB) (err error) {
	ID := uuid.New()
	c.ID = &ID
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	return nil
}
