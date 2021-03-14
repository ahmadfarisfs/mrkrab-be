package handler

//project should be nowehere here, it should be on higher level
import (
	"log"
	"net/http"
	"strconv"

	"github.com/ahmadfarisfs/krab-core/utils"
	"github.com/labstack/echo/v4"
)

//ListProject is superset for all get
func (h *Handler) ListPayRec(c echo.Context) error {
	payload := c.Get("parsedQuery").(utils.CommonRequest)
	log.Println(payload)
	res, totalData, err := h.payRecStore.ListPayRec(payload)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.Response().Header().Set("Content-Range", "payrecs "+strconv.Itoa(payload.StartIndex)+"-"+strconv.Itoa(payload.EndIndex)+"/"+strconv.Itoa(totalData))
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) Approve(c echo.Context) error {
	req := &createPayRecApproveRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	payRecID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	prd, err := h.payRecStore.Approve(uint(payRecID), int(req.SourceProjectAccountID), int(req.DestProjectAccountID), int(req.SourceBankAccountID), int(req.DestBankAccountID), &req.TrxDate)
	if err != nil {
		// log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// c.Response().Header().Set("Content-Range", "projects "+strconv.Itoa(payload.StartIndex)+"-"+strconv.Itoa(payload.EndIndex)+"/"+strconv.Itoa(totalData))
	return c.JSON(http.StatusOK, prd)
}
func (h *Handler) Reject(c echo.Context) error {
	payRecID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	prd, err := h.payRecStore.Reject(uint(payRecID))
	if err != nil {
		// log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// c.Response().Header().Set("Content-Range", "projects "+strconv.Itoa(payload.StartIndex)+"-"+strconv.Itoa(payload.EndIndex)+"/"+strconv.Itoa(totalData))
	return c.JSON(http.StatusOK, prd)
}

func (h *Handler) CreatePayRec(c echo.Context) error {
	req := &createPayRecRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	IsReimburse := false
	if req.TargetUserAccountID != nil {
		IsReimburse = true
	}
	ret, err := h.payRecStore.CreatePayRec(req.Remarks, req.Notes, req.Meta, IsReimburse, req.Amount, (req.ProjectID),
		req.SourceProjectAccountID,
		req.DestProjectAccountID,
		req.SourceBankAccountID,
		req.DestBankAccountID,
	)

	//req.BudgetID, req.SoD, req.TargetUserAccountID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, ret)
}
