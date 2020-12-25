package contract

import (
	"github.com/ahmadfarisfs/krab-core/model"
	"github.com/ahmadfarisfs/krab-core/utils"
)

type MutationStore interface {
	ListMutation(req utils.CommonRequest) ([]model.MutationExtended, int, error)
}
