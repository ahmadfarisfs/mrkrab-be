package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ahmadfarisfs/krab-core/utils"
	"github.com/labstack/echo/v4"
)

//RegisterAccount freate new accounts
func (h *Handler) RegisterFinancialAccount(c echo.Context) error {
	req := &createFinancialAccountRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	ac, err := h.accountStore.CreateAccount(req.Name, req.Description, req.AccountType, req.ParentAccount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, utils.StandardResponse{Success: true, Data: ac})
}

//ViewAccountSummary view current summary of an account
func (h *Handler) ViewFinancialAccountSummary(c echo.Context) error {
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
func (h *Handler) ListFinancialAccount(c echo.Context) error {
	payload := c.Get("parsedQuery").(utils.CommonRequest)

	res, totalData, err := h.accountStore.ListAccount(payload)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.Response().Header().Set("Access-Control-Expose-Headers", "Content-Range")
	c.Response().Header().Set("Content-Range", "projects "+strconv.Itoa(payload.StartIndex)+"-"+strconv.Itoa(payload.EndIndex)+"/"+strconv.Itoa(totalData))
	return c.JSON(http.StatusOK, res)
}

//TODO: update and delete
func (h *Handler) UpdateFinancialAccount(c echo.Context) error {

	return nil
}

func (h *Handler) DeleteFinancialAccount(c echo.Context) error {

	return nil
}
