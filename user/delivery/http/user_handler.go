package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/ahmadfarisfs/mrkrab-be/domain"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// UserHandler  represent the httphandler for User
type UserHandler struct {
	AUsecase domain.UserUsecase
}

// NewUserHandler will initialize the Users/ resources endpoint
func NewUserHandler(e *echo.Echo, us domain.UserUsecase) {
	handler := &UserHandler{
		AUsecase: us,
	}
	e.GET("/user", handler.FetchUser)
	e.POST("/user", handler.Register)
	e.GET("/user/:id", handler.GetByID)
	e.DELETE("/user/:id", handler.Delete)
}

// FetchUser will fetch the User based on given params
func (a *UserHandler) FetchUser(c echo.Context) error {
	numS := c.QueryParam("limit")
	num, _ := strconv.Atoi(numS)
	pageS := c.QueryParam("page")
	page, _ := strconv.Atoi(pageS)

	ctx := c.Request().Context()

	users, err := a.AUsecase.Fetch(ctx, int64(num), int64(page))
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	//c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, users)
}

// GetByID will get User by given id
func (a *UserHandler) GetByID(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()

	art, err := a.AUsecase.GetByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, art)
}

func isRequestValid(m *domain.User) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Register will Register the User by given request body
func (a *UserHandler) Register(c echo.Context) (err error) {
	var User domain.User
	err = c.Bind(&User)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&User); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = a.AUsecase.Register(ctx, &User)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, User)
}

// Delete will delete User by given param
func (a *UserHandler) Delete(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()

	err = a.AUsecase.Delete(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
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
