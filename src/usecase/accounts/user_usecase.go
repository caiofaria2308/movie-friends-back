package usecase_accounts

import (
	entity_accounts "app/entity/accounts"
	"fmt"
)

type userUseCase struct {
	repo IRepositoryUser
}

func NewUserUseCase(repo IRepositoryUser) IUseCaseUser {
	return &userUseCase{repo: repo}
}

func (u *userUseCase) Register(user *entity_accounts.User) error {
	err := u.repo.Create(user)
	if err != nil {
		return fmt.Errorf("could not create user")
	}
	return nil
}

func (u *userUseCase) Login(email string, password string) (*entity_accounts.User, error) {
	user, err := u.repo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}
	if err := user.CheckPassword(password); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}
	return user, nil
}

func (u *userUseCase) FindById(id int) (*entity_accounts.User, error) {
	user, err := u.repo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (u *userUseCase) FindByEmail(email string) (*entity_accounts.User, error) {
	user, err := u.repo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (u *userUseCase) Update(user *entity_accounts.User) error {
	err := u.repo.Update(user)
	if err != nil {
		return fmt.Errorf("could not update user")
	}
	return nil
}
