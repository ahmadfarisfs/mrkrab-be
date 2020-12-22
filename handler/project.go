package handler

//project should be nowehere here, it should be on higher level
import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ahmadfarisfs/krab-core/utils"
	"github.com/labstack/echo/v4"
)

//ListProject is superset for all get
func (h *Handler) ListProject(c echo.Context) error {
	payload := c.Get("parsedQuery").(utils.CommonRequest)
	log.Println(payload)
	res, totalData, err := h.projectStore.ListProject(payload)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	//	c.Response().Header().Set("Access-Control-Allow-Origin", Origin)
	//	c.Response().Header().Set("Access-Control-Allow-Methods", "GET,DELETE,POST,PUT")
	//	c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
	//	c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("Content-Range", "projects "+strconv.Itoa(payload.StartIndex)+"-"+strconv.Itoa(payload.EndIndex)+"/"+strconv.Itoa(totalData))
	//Access-Control-Expose-Headers

	return c.JSON(http.StatusOK, res)

}

func (h *Handler) GetProject(c echo.Context) error {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	//rewrite commonrequest inject using context to pass id
	res, _, err := h.projectStore.ListProject(utils.CommonRequest{Filter: map[string]interface{}{
		"id": projectID,
	}})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if len(res) == 0 {
		return c.JSON(http.StatusNotFound, nil)

	}
	return c.JSON(http.StatusOK, res[0])
}

func (h *Handler) CreateProject(c echo.Context) error {
	req := &createProjectRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	//create account for project
	account, err := h.accountStore.CreateAccount("PROJECT-"+strings.ToUpper(req.Name), nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	//create project on that account
	ac, err := h.projectStore.CreateProject(req.Name, int(account.ID), req.TotalBudget, req.Description)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	for _, p := range req.Budgets {
		accIDProjUint := uint(ac.AccountID)
		//create account for pocket
		pocketName := "PROJECT-" + strings.ToUpper(req.Name) + "-" + strings.ToUpper(p.Name)
		account, err = h.accountStore.CreateAccount(pocketName, &accIDProjUint)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		//create pocket on
		_, err := h.projectStore.CreatePocket(int(ac.ID), p.Name, account.ID, p.Budget)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	prj, _, _, err := h.projectStore.GetProjectDetails(int(ac.ID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, prj)
}

func (h *Handler) CreatePocket(c echo.Context) error {
	req := &createPocketRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	//get prjt
	project, projAccountID, _, err := h.projectStore.GetProjectDetails(req.ProjectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	//create account for project
	pocketName := "PROJECT-" + strings.ToUpper(project.Name) + "-" + strings.ToUpper(req.Name)
	account, err := h.accountStore.CreateAccount(pocketName, &projAccountID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	//create pocket on
	prj, err := h.projectStore.CreatePocket(req.ProjectID, req.Name, account.ID, req.Budget)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.StandardResponse{Success: true, Data: prj})
}

func (h *Handler) CreateProjectTransaction(c echo.Context) error {
	req := &createProjectTransactionRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	//get project info
	_, projAccountID, budgetAccountIDs, err := h.projectStore.GetProjectDetails(req.ProjectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var accountID int
	if req.BudgetID != nil {
		//check that account should be under projects
		isValid := false
		for _, v := range budgetAccountIDs {
			if v == *req.BudgetID {
				//good
				isValid = true
			}
		}
		if !isValid {
			return c.JSON(http.StatusInternalServerError, utils.StandardResponse{Success: false, ErrorMessage: "invalid budgetid"})
		}
		accountID = int(*req.BudgetID)
	} else {
		accountID = int(projAccountID)
	}

	trx, err := h.transactionStore.CreateTransaction(accountID, req.Amount, req.Remarks)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.StandardResponse{Success: true, Data: trx})
}

func (h *Handler) CreateProjectTransfer(c echo.Context) error {
	//todo: refactor this, bussiness logic should be nowhere here

	req := &createProjectTransferRequest{}
	var sourceAccount int
	var targetAccount int
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	trfDir, isSameProject := req.analyze()
	if !isSameProject && trfDir != ProjectToProject {
		//illegal
		return c.JSON(http.StatusUnauthorized, utils.StandardResponse{Success: false, ErrorMessage: "Inter project can only transfered to project account"})
	}

	//get project and its accounts
	_, projectSourceAccountID, _, err := h.projectStore.GetProjectDetails(req.ProjectIDSource)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}
	_, projectTargetAccountID, _, err := h.projectStore.GetProjectDetails(req.ProjectIDTarget)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	//sof: project
	if trfDir == ProjectToPocket || trfDir == ProjectToProject {
		_, _, _, err := h.projectStore.GetProjectDetails(req.ProjectIDSource)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, utils.StandardResponse{Success: false, ErrorMessage: "Invalid Source Project ID"})
		}
		sourceAccount = int(projectSourceAccountID)
	}

	//sof: pocket
	if trfDir == PocketToProject || trfDir == PocketToPocket {
		//check budgetid validity
		_, err := h.projectStore.CheckBudgetIDValidity(int(*req.BudgetIDSource), req.ProjectIDSource)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, utils.StandardResponse{Success: false, ErrorMessage: "Invalid Source Budget ID"})
		}
		sourceAccount = int(*req.BudgetIDSource)
	}

	//tof: project
	if trfDir == PocketToProject || trfDir == ProjectToProject {
		_, _, _, err := h.projectStore.GetProjectDetails(req.ProjectIDTarget)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, utils.StandardResponse{Success: false, ErrorMessage: "Invalid Target Project ID"})
		}
		targetAccount = int(projectTargetAccountID)
	}

	//tof: pocket
	if trfDir == ProjectToPocket || trfDir == PocketToPocket {
		_, err := h.projectStore.CheckBudgetIDValidity(int(*req.BudgetIDSource), req.ProjectIDSource)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, utils.StandardResponse{Success: false, ErrorMessage: "Invalid Target Budget ID"})
		}

		targetAccount = int(*req.BudgetIDTarget)
	}

	ret, err := h.transactionStore.CreateTransfer(sourceAccount, targetAccount, req.Amount, req.Remarks)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, utils.StandardResponse{Success: true, Data: ret})
}
