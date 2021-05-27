package handler

import (
	"time"

	"github.com/ahmadfarisfs/krab-core/contract"
	"github.com/labstack/echo/v4"
)

type createUserRequest struct {
	Fullname string `validate:"required"`
	Username string `validate:"required"`
	Role     string `validate:"required"`
	Password string `validate:"required"`
	Email    string `validate:"required"`
}

func (ca *createUserRequest) bind(c echo.Context) error {
	if err := c.Bind(ca); err != nil {
		return err
	}
	if err := c.Validate(ca); err != nil {
		return err
	}
	return nil
}

type createFinancialAccountRequest struct {
	Name          string `validate:"required"`
	ParentAccount *uint  `validate:"omitempty"`
	AccountType   string `validate:"oneof=internal external"`
	Description   string `validate:"required"`
}

func (ca *createFinancialAccountRequest) bind(c echo.Context) error {
	if err := c.Bind(ca); err != nil {
		return err
	}
	if err := c.Validate(ca); err != nil {
		return err
	}
	return nil
}

type createIncomeTransactionRequest struct {
	Amount                   int       `validate:"required"`
	Remarks                  string    `validate:"required"`
	DestinationProjectID     int       `validate:"required"`
	IncomeFinancialAccountID int       `validate:"required"`
	SourceBankAccountID      int       `validate:"required"`
	DestinationBankAccountID int       `validate:"required"`
	TransactionTime          time.Time `validate:"required"`
}

func (ca *createIncomeTransactionRequest) bind(c echo.Context) error {
	if err := c.Bind(ca); err != nil {
		return err
	}
	if err := c.Validate(ca); err != nil {
		return err
	}
	return nil
}

type createExpenseTransactionRequest struct {
	Amount              int                           `validate:"required"`
	Remarks             string                        `validate:"required"`
	SourceProjectID     int                           `validate:"required"`
	SourceBankAccountID int                           `validate:"required"`
	Expenses            []contract.ExpenseDestination `validate:"required"`
	TransactionTime     time.Time                     `validate:"required"`
}

func (ca *createExpenseTransactionRequest) bind(c echo.Context) error {
	if err := c.Bind(ca); err != nil {
		return err
	}
	if err := c.Validate(ca); err != nil {
		return err
	}
	return nil
}

type createBankTransferTransactionRequest struct {
	Amount               int                            `validate:"required"`
	Remarks              string                         `validate:"required"`
	TransferFeeProjectID int                            `validate:"required"`
	SourceBankAccountID  int                            `validate:"required"`
	Destination          []contract.TransferDestination `validate:"required"`
	TransferFee          contract.ExpenseDestination    `validate:"required"`
	TransactionTime      time.Time                      `validate:"required"`
}

func (ca *createBankTransferTransactionRequest) bind(c echo.Context) error {
	if err := c.Bind(ca); err != nil {
		return err
	}
	if err := c.Validate(ca); err != nil {
		return err
	}
	return nil
}

type createProjectTransferTransactionRequest struct {
	Amount                        int       `validate:"required"`
	Remarks                       string    `validate:"required"`
	SourceProjectID               int       `validate:"required"`
	SourceFinancialAccountID      int       `validate:"required"`
	DestinationProjectID          int       `validate:"required"`
	DestinationFinancialAccountID int       `validate:"required"`
	TransactionTime               time.Time `validate:"required"`
}

func (ca *createProjectTransferTransactionRequest) bind(c echo.Context) error {
	if err := c.Bind(ca); err != nil {
		return err
	}
	if err := c.Validate(ca); err != nil {
		return err
	}
	return nil
}

type createTransferRequest struct {
	AccountFrom     int       `validate:"required"`
	AccountTo       int       `validate:"required"`
	Amount          uint      `validate:"required"`
	Remarks         string    `validate:"required"`
	TransactionDate time.Time `validate:"required"`
}

func (ca *createTransferRequest) bind(c echo.Context) error {
	if err := c.Bind(ca); err != nil {
		return err
	}
	if err := c.Validate(ca); err != nil {
		return err
	}
	return nil
}

type createProjectRequest struct {
	TotalBudget int    `validate:"required"`
	Name        string `validate:"required"`
	Description string `validate:"required"`
	Status      string `validate:"required"`
	Budgets     []createPocketRequest
}

func (ca *createProjectRequest) bind(c echo.Context) error {
	if err := c.Bind(ca); err != nil {
		return err
	}
	if err := c.Validate(ca); err != nil {
		return err
	}
	return nil
}

type createPocketRequest struct {
	ProjectID int    `validate:"omitempty"`
	Budget    *uint  `validate:"omitempty"`
	Name      string `validate:"required"`
}

func (ca *createPocketRequest) bind(c echo.Context) error {
	if err := c.Bind(ca); err != nil {
		return err
	}
	if err := c.Validate(ca); err != nil {
		return err
	}
	return nil
}

type updateProjectRequest struct {
	ProjectID   int    `validate:"required"`
	Status      string `validate:"required"`
	Description string `validate:"required"`
	TotalBudget int    `validate:"required"`
	Name        string `validate:"required"`
}

func (ca *updateProjectRequest) bind(c echo.Context) error {
	if err := c.Bind(ca); err != nil {
		return err
	}
	if err := c.Validate(ca); err != nil {
		return err
	}
	return nil
}

type createProjectTransactionRequest struct {
	ProjectID       int       `validate:"required"`
	BudgetID        *uint     `validate:"omitempty"`
	Amount          int       `validate:"required"`
	Remarks         string    `validate:"required"`
	URL             string    `validate:"omitempty"`
	Notes           string    `validate:"omitempty"`
	Meta            string    `validate:"omitempty"`
	SoD             string    `validate:"required"`
	TransactionDate time.Time `validate:"required"`
}

func (ca *createProjectTransactionRequest) bind(c echo.Context) error {
	if err := c.Bind(ca); err != nil {
		return err
	}
	if err := c.Validate(ca); err != nil {
		return err
	}
	return nil
}

type createProjectTransferRequest struct {
	ProjectIDSource int   `validate:"required"`
	BudgetIDSource  *uint `validate:"omitempty"`
	ProjectIDTarget int   `validate:"required"`
	BudgetIDTarget  *uint `validate:"omitempty"`

	Amount  uint      `validate:"required"`
	Remarks string    `validate:"required"`
	URL     string    `validate:"omitempty"`
	Notes   string    `validate:"omitempty"`
	Meta    string    `validate:"omitempty"`
	TrxDate time.Time `validate:"required"`
}

func (ca *createProjectTransferRequest) bind(c echo.Context) error {
	if err := c.Bind(ca); err != nil {
		return err
	}
	if err := c.Validate(ca); err != nil {
		return err
	}
	return nil
}

//budget can only transfered within the same projects
type TransferDirection int

const (
	PocketToPocket TransferDirection = iota
	ProjectToPocket
	PocketToProject
	ProjectToProject
	Invalid
)

func (ca *createProjectTransferRequest) analyze() (dir TransferDirection, isSameProject bool) {
	if ca.ProjectIDSource == ca.ProjectIDTarget {
		isSameProject = true
	}
	if ca.BudgetIDSource == nil && ca.BudgetIDTarget == nil {
		dir = ProjectToProject
	} else if ca.BudgetIDSource != nil && ca.BudgetIDTarget == nil {
		dir = PocketToProject
	} else if ca.BudgetIDSource != nil && ca.BudgetIDTarget != nil {
		dir = PocketToPocket
	} else if ca.BudgetIDSource == nil && ca.BudgetIDTarget != nil {
		dir = ProjectToPocket
	} else {
		dir = Invalid
	}

	return

}

type createPayRecRequest struct {
	ProjectID int    `validate:"required"`
	BudgetID  *uint  `validate:"omitempty"`
	Remarks   string `validate:"required"`
	Amount    int    `validate:"required"`
	SoD       string `validate:"required"`
}

func (ca *createPayRecRequest) bind(c echo.Context) error {
	if err := c.Bind(ca); err != nil {
		return err
	}
	if err := c.Validate(ca); err != nil {
		return err
	}
	return nil
}

type loginRequest struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

func (ca *loginRequest) bind(c echo.Context) error {
	if err := c.Bind(ca); err != nil {
		return err
	}
	if err := c.Validate(ca); err != nil {
		return err
	}
	return nil
}

type createBankAccountRequest struct {
	BankName      string `validate:"required"`
	HolderName    string `validate:"required"`
	AccountNumber string `validate:"required"`
	Description   string `validate:"required"`
	AccountType   string `validate:"required"`
}

func (ca *createBankAccountRequest) bind(c echo.Context) error {
	if err := c.Bind(ca); err != nil {
		return err
	}
	if err := c.Validate(ca); err != nil {
		return err
	}
	return nil
}
