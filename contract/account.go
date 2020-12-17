package contract

import (
	"github.com/ahmadfarisfs/krab-core/model"
	"github.com/ahmadfarisfs/krab-core/utils"
)

type AccountStore interface {
	ListAccount(req utils.CommonRequest) ([]model.Account, int, error)
	CreateAccount(name string, parentID *uint) (model.Account, error)
	GetAccountDetails(id int) (model.Account, error)
}
