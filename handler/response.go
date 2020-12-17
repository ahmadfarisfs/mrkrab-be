package handler

import "time"

type pagingResponse struct {
	TotalData int         `json:"total"`
	Data      interface{} `json:"data"`
}

type createTransactionResponse struct {
	ID              uint
	TransactionCode string
	FromID          int
	FromName        string
	ToID            int
	ToName          string
	Amount          int
	Remarks         string
}
type detailsProjectResponse struct {
	ID                         uint
	Name                       string
	IsOpen                     bool
	Amount                     *uint
	Description                *string
	ProjectAccountBalance      int
	ProjectAccountTotalIncome  int
	ProjectAccountTotalExpense int
	Budgets                    []detailsBudgetReponse
	CreatedOn                  time.Time
}

type detailsBudgetReponse struct {
	ID           int
	Name         string
	Limit        *uint
	Balance      int
	TotalIncome  int
	TotalExpense int
	CreatedOn    time.Time
}
