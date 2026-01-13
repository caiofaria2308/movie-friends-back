package entity_crew

import (
	entity_accounts "app/entity/accounts"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Crew struct {
	ID        *uuid.UUID            `json:"id"`
	Name      string                `json:"name"`
	Owner     *entity_accounts.User `json:"owner"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

func (c *Crew) TableName() string {
	return "crew_crews"
}

func (c *Crew) BeforeCreate(tx *gorm.DB) (err error) {
	ID := uuid.New()
	c.ID = &ID
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}
