package entity_accounts

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserPix struct {
	ID        *uuid.UUID `json:"id"`
	Owner     *User      `json:"owner"`
	OwnerID   int        `json:"owner_id"`
	PixKey    string     `json:"pix_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (c *UserPix) TableName() string {
	return "accounts_user_pix"
}

func (c *UserPix) BeforeCreate(tx *gorm.DB) (err error) {
	ID := uuid.New()
	c.ID = &ID
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}
