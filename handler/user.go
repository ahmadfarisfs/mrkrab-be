package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ahmadfarisfs/krab-core/utils"
	"github.com/labstack/echo/v4"
)

//ListProject is superset for all get
func (h *Handler) ListUser(c echo.Context) error {
	payload := c.Get("parsedQuery").(utils.CommonRequest)
	log.Println(payload)
	res, totalData, err := h.userStore.ListUser(payload)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	//	c.Response().Header().Set("Access-Control-Allow-Origin", Origin)
	//	c.Response().Header().Set("Access-Control-Allow-Methods", "GET,DELETE,POST,PUT")
	//	c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
	//	c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("Content-Range", "users "+strconv.Itoa(payload.StartIndex)+"-"+strconv.Itoa(payload.EndIndex)+"/"+strconv.Itoa(totalData))
	//Access-Control-Expose-Headers

	return c.JSON(http.StatusOK, res)

}

func (h *Handler) GetUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	//rewrite commonrequest inject using context to pass id
	res, _, err := h.userStore.ListUser(utils.CommonRequest{Filter: map[string]interface{}{
		"id": userID,
	}})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if len(res) == 0 {
		return c.JSON(http.StatusNotFound, nil)

	}
	return c.JSON(http.StatusOK, res[0])
}

//CreateTransaction
func (h *Handler) CreateUser(c echo.Context) error {
	req := &createUserRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	//name string, username string, password string, email string, role string
	ac, err := h.userStore.CreateUser(req.Fullname, req.Username, req.Password, req.Email, req.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ac)
}

func (h *Handler) DeleteUser(c echo.Context) error {
	//TODO: check user validity
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	err = h.userStore.DeleteUser(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}
