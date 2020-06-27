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

// ProjectHandler  represent the httphandler for Project
type ProjectHandler struct {
	AUsecase domain.ProjectUsecase
}

// NewProjectHandler will initialize the Projects/ resources endpoint
func NewProjectHandler(e *echo.Echo, us domain.ProjectUsecase) {
	handler := &ProjectHandler{
		AUsecase: us,
	}
	e.GET("/project", handler.FetchProject)
	e.POST("/project", handler.Register)
	e.GET("/project/{id:[0-9]+}", handler.GetByID)
	e.DELETE("/project/:id", handler.Delete)
}

// FetchProject will fetch the Project based on given params
func (a *ProjectHandler) FetchProject(c echo.Context) error {
	numS := c.QueryParam("limit")
	num, _ := strconv.Atoi(numS)
	pageS := c.QueryParam("page")
	page, _ := strconv.Atoi(pageS)

	ctx := c.Request().Context()

	projects, totalRecord, totalPage, err := a.AUsecase.Fetch(ctx, int64(num), int64(page), nil)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	//c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, utilities.Paginator{
		TotalPage:   totalPage,
		TotalRecord: totalRecord,
		Records:     projects,
	})
}

// GetByID will get Project by given id
func (a *ProjectHandler) GetByID(c echo.Context) error {
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

func isRequestValid(m *domain.Project) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Register will Register the Project by given request body
func (a *ProjectHandler) Register(c echo.Context) (err error) {
	var Project domain.Project

	err = c.Bind(&Project)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&Project); !ok {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = a.AUsecase.Add(ctx, &Project)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	log.Println(Project)
	return c.JSON(http.StatusCreated, Project)
}

// Delete will delete Project by given param
func (a *ProjectHandler) Delete(c echo.Context) error {
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
