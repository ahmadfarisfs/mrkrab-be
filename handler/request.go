package handler

import (
	"time"

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

type createProjectAccountRequest struct {
	Name        string `validate:"required"`
	AccountType string `validate:"oneof=EXPENSE INCOME"` //expense or income
	Meta        string `validate:"omitempty"`
}

func (ca *createProjectAccountRequest) bind(c echo.Context) error {
	if err := c.Bind(ca); err != nil {
		return err
	}
	if err := c.Validate(ca); err != nil {
		return err
	}
	return nil
}

type createBankAccountRequest struct {
	BankName       string `validate:"required"`
	BankNumber     string `validate:"required"`
	BankHoldername string `validate:"required"`
	Internal       bool
	Meta           string `validate:"omitempty"`
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

type createAccountRequest struct {
	Name          string `validate:"required"`
	ParentAccount *uint  `validate:"omitempty"`
	Meta          string `validate:"omitempty"`
	AccountType   string `validate:"oneof=BANK EXPENSE PROJECT INCOME"`
}

func (ca *createAccountRequest) bind(c echo.Context) error {
	if err := c.Bind(ca); err != nil {
		return err
	}
	if err := c.Validate(ca); err != nil {
		return err
	}
	return nil
}

type createTransactionRequest struct {
	AccountID       int       `validate:"required"`
	Amount          int       `validate:"required"`
	Remarks         string    `validate:"required"`
	SoD             string    `validate:"required"`
	TransactionDate time.Time `validate:"required"`
}

func (ca *createTransactionRequest) bind(c echo.Context) error {
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
	// TotalBudget *uint   `validate:"omitempty"`
	Name        string `validate:"required"`
	Description string `validate:"required"`
	Amount      *uint  `validate:"omitempty"`
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
	ProjectID       int  `validate:"required"`
	BudgetAccountID int  `validate:"required"`
	Amount          uint `validate:"required"`
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
	Description string `validate:"omitempty"`
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
	ProjectID int `validate:"required"`
	Amount    int `validate:"required"`

	SourceProjectAccountID int `validate:"required"`
	DestProjectAccountID   int `validate:"required"`

	SourceBankAccountID int `validate:"required"`
	DestBankAccountID   int `validate:"required"`

	Remarks         string    `validate:"required"`
	Notes           string    `validate:"omitempty"`
	Meta            string    `validate:"omitempty"`
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

type createPayRecApproveRequest struct {
	SourceBankAccountID uint `validate:"required"`
	DestBankAccountID   uint `validate:"required"`

	SourceProjectAccountID uint `validate:"required"`
	DestProjectAccountID   uint `validate:"required"`

	Remarks string    `validate:"required"`
	Notes   string    `validate:"omitempty"`
	Meta    string    `validate:"omitempty"`
	TrxDate time.Time `validate:"required"`
}

func (ca *createPayRecApproveRequest) bind(c echo.Context) error {
	if err := c.Bind(ca); err != nil {
		return err
	}
	if err := c.Validate(ca); err != nil {
		return err
	}
	return nil
}

type createBankTransferRequest struct {
	SourceBankAccountID uint `validate:"required"`
	DestBankAccountID   uint `validate:"required"`

	Amount  uint      `validate:"required"`
	Remarks string    `validate:"required"`
	Notes   string    `validate:"omitempty"`
	Meta    string    `validate:"omitempty"`
	TrxDate time.Time `validate:"required"`
}

func (ca *createBankTransferRequest) bind(c echo.Context) error {
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
	ProjectID int `validate:"required"`

	SourceBankAccountID    *int `validate:"omitempty"`
	DestBankAccountID      *int `validate:"omitempty"`
	SourceProjectAccountID *int `validate:"omitempty"`
	DestProjectAccountID   *int `validate:"omitempty"`

	Remarks string `validate:"required"`
	Amount  int    `validate:"required"`
	Notes   string `validate:"omitempty"`
	Meta    string `validate:"omitempty"`

	TargetUserAccountID *int `validate:"omitempty"`
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
