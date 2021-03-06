package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ahmadfarisfs/krab-core/utils"
	"github.com/labstack/echo/v4"
)

//ListProject is superset for all get
func (h *Handler) ListFinancialMutation(c echo.Context) error {
	payload := c.Get("parsedQuery").(utils.CommonRequest)
	log.Println(payload)
	res, totalData, err := h.mutationStore.ListFinancialAccountMutation(payload)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.Response().Header().Set("Content-Range", "users "+strconv.Itoa(payload.StartIndex)+"-"+strconv.Itoa(payload.EndIndex)+"/"+strconv.Itoa(totalData))
	return c.JSON(http.StatusOK, res)

}

func (h *Handler) ListBankMutation(c echo.Context) error {
	payload := c.Get("parsedQuery").(utils.CommonRequest)
	log.Println(payload)
	res, totalData, err := h.mutationStore.ListBankAccountMutation(payload)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.Response().Header().Set("Content-Range", "users "+strconv.Itoa(payload.StartIndex)+"-"+strconv.Itoa(payload.EndIndex)+"/"+strconv.Itoa(totalData))
	return c.JSON(http.StatusOK, res)

}
