package handler

import "github.com/labstack/echo/v4"

type createAccountRequest struct {
	Name          string `validate:"required"`
	ParentAccount *uint  `validate:"omitempty"`
	//	AccountType string `validate:"oneof=assets expenses liabilities revenues"`
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
	AccountID int    `validate:"required"`
	Amount    int    `validate:"required"`
	Remarks   string `validate:"required"`
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
	AccountFrom int    `validate:"required"`
	AccountTo   int    `validate:"required"`
	Amount      uint   `validate:"required"`
	Remarks     string `validate:"required"`
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
	Budget      *uint   `validate:"omitempty"`
	Name        string  `validate:"required"`
	Description *string `validate:"omitempty"`
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

type createProjectTransactionRequest struct {
	ProjectID int    `validate:"required"`
	BudgetID  *uint  `validate:"omitempty"`
	Amount    int    `validate:"required"`
	Remarks   string `validate:"required"`
	URL       string `validate:"omitempty"`
	Notes     string `validate:"omitempty"`
	Meta      string `validate:"omitempty"`
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

	Amount  uint   `validate:"required"`
	Remarks string `validate:"required"`
	URL     string `validate:"omitempty"`
	Notes   string `validate:"omitempty"`
	Meta    string `validate:"omitempty"`
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
