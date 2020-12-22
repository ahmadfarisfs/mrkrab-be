package contract

import (
	"github.com/ahmadfarisfs/krab-core/model"
	"github.com/ahmadfarisfs/krab-core/utils"
)

type UserStore interface {
	ListUser(req utils.CommonRequest) ([]model.User, int, error)
	CreateUser(name string, username string, password string, email string, role string) (model.User, error)
	//GetUserDetails(id int) (model.User, error)
	DeleteUser(id int) error
}
