package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ahmadfarisfs/krab-core/utils"
	"github.com/labstack/echo/v4"
)

//RegisterAccount freate new accounts
func (h *Handler) RegisterAccount(c echo.Context) error {
	req := &createAccountRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	ac, err := h.accountStore.CreateAccount(req.Name, req.ParentAccount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, utils.StandardResponse{Success: true, Data: ac})
}

//ViewAccountSummary view current summary of an account
func (h *Handler) ViewAccountSummary(c echo.Context) error {
	accountID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	ac, err := h.accountStore.GetAccountDetails(accountID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, utils.StandardResponse{Success: true, Data: ac})
}

//ListAccount list account with paging
func (h *Handler) ListAccount(c echo.Context) error {
	payload := c.Get("parsedQuery").(utils.CommonRequest)

	res, totalData, err := h.accountStore.ListAccount(payload)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.Response().Header().Set("Access-Control-Expose-Headers", "Content-Range")
	c.Response().Header().Set("Content-Range", "projects "+strconv.Itoa(payload.StartIndex)+"-"+strconv.Itoa(payload.EndIndex)+"/"+strconv.Itoa(totalData))
	// /c.Response().Header().Set("Access-Control-Expose-Headers", "Content-Range")

	return c.JSON(http.StatusOK, res)
}
