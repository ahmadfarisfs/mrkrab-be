package contract

import (
	"github.com/ahmadfarisfs/krab-core/model"
	"github.com/ahmadfarisfs/krab-core/utils"
)

type PayRecStore interface {
	CreatePayRec(remarks string, amount int, projectID uint, pocketID *uint) (model.PayRec, error)
	Approve(id uint) (model.PayRec, error)
	Reject(id uint) (model.PayRec, error)
	ListPayRec(req utils.CommonRequest) ([]model.PayRec, int, error)
}
