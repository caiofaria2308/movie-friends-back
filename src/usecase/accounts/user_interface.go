package usecase_accounts

import entity_accounts "app/entity/accounts"

type IRepositoryUser interface {
	Create(user *entity_accounts.User) error
	FindById(id int) (*entity_accounts.User, error)
	FindByEmail(email string) (*entity_accounts.User, error)
	Update(user *entity_accounts.User) error
}

type IUseCaseUser interface {
	Register(user *entity_accounts.User) error
	Login(email string, password string) (*entity_accounts.User, error)
	FindById(id int) (*entity_accounts.User, error)
	FindByEmail(email string) (*entity_accounts.User, error)
	Update(user *entity_accounts.User) error
}
