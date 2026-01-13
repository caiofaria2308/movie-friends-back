package entity_accounts

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	ROLE_ADMIN = "admin"
	ROLE_USER  = "user"
	ROLE_GUEST = "guest"
)

type User struct {
	ID        int       `gorm:"primarykey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Role      string    `json:"role"`
}

func (User) TableName() string {
	return "account_users"
}

func (u *User) EncryptedPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
