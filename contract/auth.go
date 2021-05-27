package contract

import "github.com/ahmadfarisfs/krab-core/model"

type AuthStore interface {
	Login(username, password string) error
	CreateUser(name string) (model.User, error)
	UpdateUser(id int, name string, role string) error
	GetUserDetails(id int) (model.User, error)
}
