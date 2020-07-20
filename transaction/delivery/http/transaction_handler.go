package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/ahmadfarisfs/mrkrab-be/domain"
	"github.com/ahmadfarisfs/mrkrab-be/utilities"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// TransactionHandler  represent the httphandler for User
type TransactionHandler struct {
	TrxUsecase domain.TransactionUsecase
}

// NewTransactionHandler will initialize the Users/ resources endpoint
func NewTransactionHandler(e *echo.Echo, us domain.TransactionUsecase) {
	handler := &TransactionHandler{
		TrxUsecase: us,
	}
	e.POST("/transaction", handler.AddTransaction)
	e.GET("/transaction/:id", handler.GetByID)
	//	e.GET("/transaction", handler.Register)
}

// FetchTransaction will fetch the User based on given params
func (a *TransactionHandler) FetchTransaction(c echo.Context) error {
	numS := c.QueryParam("limit")
	num, _ := strconv.Atoi(numS)
	pageS := c.QueryParam("page")
	page, _ := strconv.Atoi(pageS)
	filter := domain.Transaction{}

	projectIDS := c.QueryParam("project_id")
	if projectIDS != "" {
		projectID, err := strconv.Atoi(projectIDS)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		filter.ProjectID = projectID
	}

	createdBy := c.QueryParam("created_by")
	if createdBy != "" {
		createdBy, err := strconv.Atoi(createdBy)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		filter.CreatorID = createdBy
	}

	ctx := c.Request().Context()
	trx, totalRecord, totalPage, err := a.TrxUsecase.Fetch(ctx, int64(num), int64(page), &filter)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, utilities.Paginator{
		TotalPage:   totalPage,
		TotalRecord: totalRecord,
		Records:     trx,
	})
}

// AddTransaction will Add the Project by given request body
func (a *TransactionHandler) AddTransaction(c echo.Context) (err error) {
	var Transaction domain.Transaction

	//err = utilities.BindAndValidate(c, Transaction)
	log.Println("here")
	err = c.Bind(&Transaction)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	err = c.Validate(Transaction)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	/*	if ok, err := isRequestValid(&Transaction); !ok {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
	*/
	ctx := c.Request().Context()
	err = a.TrxUsecase.Add(ctx, &Transaction)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	log.Println(Transaction)
	return c.JSON(http.StatusCreated, Transaction)
}

func (a *TransactionHandler) AddBudget(c echo.Context) (err error) {
	payload := struct {
		ProjectID  int `json:"project_id" validator:"required"`
		CategoryID int `json:"category_id" validator:"required"`
		Amount     int `json:"amount" validator:"required,min=0"`
	}{}
	err = c.Bind(&payload)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	return nil
}

// GetByID will get User by given id
func (a *TransactionHandler) GetByID(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()

	art, err := a.TrxUsecase.GetByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, art)
}

func isRequestValid(m *domain.Transaction) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
