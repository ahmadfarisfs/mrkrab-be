package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ahmadfarisfs/krab-core/utils"
	"github.com/labstack/echo/v4"
)

type ListTransactionRequest struct {
}

//ViewTransactionDetails see details transaction by trx_id
func (h *Handler) ViewTransactionDetails(c echo.Context) error {
	trxID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	trx, err := h.transactionStore.GetTransactionDetailsbyID(trxID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, trx)

}

//ViewTransaction list mutation between dates with filter (list of account id) also with paging
func (h *Handler) ListTransaction(c echo.Context) error {
	payload := c.Get("parsedQuery").(utils.CommonRequest)

	if payload.Filter["start_time"] != nil {
		_, err := time.ParseInLocation(time.RFC3339, payload.Filter["start_time"].(string), utils.TimeLocation)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, err.Error())
		}
	}
	if payload.Filter["end_time"] != nil {
		_, err := time.ParseInLocation(time.RFC3339, payload.Filter["end_time"].(string), utils.TimeLocation)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, err.Error())
		}
	}

	res, totalData, err := h.transactionStore.ListTransaction(payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.Response().Header().Set("Content-Range", "projects "+strconv.Itoa(payload.StartIndex)+"-"+strconv.Itoa(payload.EndIndex)+"/"+strconv.Itoa(totalData))

	return c.JSON(http.StatusOK, res)
}

//CreateTransaction
func (h *Handler) CreateTransaction(c echo.Context) error {
	req := &createTransactionRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	ac, err := h.transactionStore.CreateTransaction(req.AccountID, req.Amount, req.Remarks, req.SoD, req.TransactionDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ac)
}

//CreateTransaction
func (h *Handler) CreateTransfer(c echo.Context) error {
	req := &createTransferRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	ac, err := h.transactionStore.CreateTransfer(req.AccountFrom, req.AccountTo, req.Amount, req.Remarks)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ac)
}
