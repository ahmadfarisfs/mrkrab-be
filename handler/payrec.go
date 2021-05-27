package handler

//project should be nowehere here, it should be on higher level
import (
	"github.com/labstack/echo/v4"
)

//ListProject is superset for all get
func (h *Handler) ListPayRec(c echo.Context) error {
	// payload := c.Get("parsedQuery").(utils.CommonRequest)
	// log.Println(payload)
	// res, totalData, err := h.payRecStore.ListPayRec(payload)
	// if err != nil {
	// 	log.Println(err)
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }
	// c.Response().Header().Set("Content-Range", "payrecs "+strconv.Itoa(payload.StartIndex)+"-"+strconv.Itoa(payload.EndIndex)+"/"+strconv.Itoa(totalData))
	// return c.JSON(http.StatusOK, res)
	return nil
}

func (h *Handler) Approve(c echo.Context) error {
	// payRecID, err := strconv.Atoi(c.Param("id"))
	// if err != nil {
	// 	return c.JSON(http.StatusUnprocessableEntity, err.Error())
	// }
	// prd, err := h.payRecStore.Approve(uint(payRecID))
	// if err != nil {
	// 	// log.Println(err)
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }
	// // c.Response().Header().Set("Content-Range", "projects "+strconv.Itoa(payload.StartIndex)+"-"+strconv.Itoa(payload.EndIndex)+"/"+strconv.Itoa(totalData))
	// return c.JSON(http.StatusOK, prd)

	return nil
}
func (h *Handler) Reject(c echo.Context) error {
	// payRecID, err := strconv.Atoi(c.Param("id"))
	// if err != nil {
	// 	return c.JSON(http.StatusUnprocessableEntity, err.Error())
	// }
	// prd, err := h.payRecStore.Reject(uint(payRecID))
	// if err != nil {
	// 	// log.Println(err)
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }
	// // c.Response().Header().Set("Content-Range", "projects "+strconv.Itoa(payload.StartIndex)+"-"+strconv.Itoa(payload.EndIndex)+"/"+strconv.Itoa(totalData))
	// return c.JSON(http.StatusOK, prd)
	return nil
}

func (h *Handler) CreatePayRec(c echo.Context) error {
	// req := &createPayRecRequest{}
	// if err := req.bind(c); err != nil {
	// 	return c.JSON(http.StatusUnprocessableEntity, err.Error())
	// }

	// ret, err := h.payRecStore.CreatePayRec(req.Remarks, req.Amount, uint(req.ProjectID), req.BudgetID, req.SoD)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }

	// return c.JSON(http.StatusOK, ret)
	return nil
}
