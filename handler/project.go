package handler

//project should be nowehere here, it should be on higher level
import (
	"log"
	"net/http"
	"strconv"

	"github.com/ahmadfarisfs/krab-core/model"
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

func (h *Handler) DeleteProject(c echo.Context) error {
	//TODO: check user validity
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	err = h.projectStore.DeleteProject(projectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
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

	budgetAccounts := []uint{}
	budgetLimits := []uint{}

	for _, v := range req.Budgets {
		budgetAccounts = append(budgetAccounts, uint(v.BudgetAccountID))
		budgetLimits = append(budgetLimits, v.Amount)
	}

	//create project with all accounts inside
	ac, err := h.projectStore.CreateProject(req.Name, req.Amount, req.Description, budgetAccounts, budgetLimits)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// accIDProjUint := uint(ac.AccountID)

	// generalAccountName := "PROJECT-" + strings.ToUpper(req.Name) + "-GENERAL-" + strconv.Itoa(int(time.Now().Unix()))

	//create default expense account
	// account, err = h.accountStore.CreateAccount(generalAccountName, &accIDProjUint)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }
	// //create pocket on general
	// _, err = h.projectStore.CreatePocket(int(ac.ID), "General", account.ID, req.TotalBudget)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }

	// //create default income account
	// incomeAccountName := "PROJECT-" + strings.ToUpper(req.Name) + "-GENERAL-" + strconv.Itoa(int(time.Now().Unix()))

	// account, err = h.accountStore.CreateAccount(incomeAccountName, &accIDProjUint)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }
	// //create pocket on general
	// _, err = h.projectStore.CreatePocket(int(ac.ID), "General", account.ID, req.TotalBudget)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }

	//create other accounts
	// for _, p := range req.Budgets {

	// 	//create account for pocket
	// 	// pocketName := "PROJECT-" + strings.ToUpper(req.Name) + "-" + strings.ToUpper(p.Name) + "-" + strconv.Itoa(int(time.Now().Unix()))
	// 	// account, err = h.accountStore.CreateAccount(pocketName, &accIDProjUint)
	// 	// if err != nil {
	// 	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// 	// }

	// 	//create pocket on
	// 	_, err := h.projectStore.CreatePocket(nil, int(ac.ID), p.Name, account.ID, p.Budget)
	// 	if err != nil {
	// 		return c.JSON(http.StatusInternalServerError, err.Error())
	// 	}
	// }

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
	_, _, _, err := h.projectStore.GetProjectDetails(req.ProjectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	//create account for project
	// pocketName := "PROJECT-" + strings.ToUpper(project.Name) + "-" + strings.ToUpper(req.Name) + "-" + strconv.Itoa(int(time.Now().Unix()))
	// account, err := h.accountStore.CreateAccount(pocketName, &projAccountID)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }

	//create pocket on
	prj, err := h.projectStore.CreatePocket(req.ProjectID, uint(req.BudgetAccountID), req.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, prj)
}

func (h *Handler) CreateProjectTransaction(c echo.Context) error {
	req := &createProjectTransactionRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	//get project info
	// proj, projAccountID, budgetAccountIDs, err := h.projectStore.GetProjectDetails(req.ProjectID)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }

	trx, err := h.transactionStore.CreateProjectTransaction(uint(req.SourceProjectAccountID), uint(req.DestProjectAccountID),
		uint(req.SourceBankAccountID), uint(req.DestBankAccountID), uint(req.Amount), req.Remarks, req.Meta, req.Notes, req.TransactionDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// if req.Amount > 0 {
	// 	//income can only come from REVENUE ACCOUNT with prefix name: ACCOUNT-REVENUE and id 0
	// 	// to projects
	// 	//find revenue account
	// 	trx, err := h.transactionStore.CreateTransfer(0, int(projAccountID), uint(math.Abs(float64(req.Amount))), req.Remarks, req.TransactionDate, false)
	// 	if err != nil {
	// 		return c.JSON(http.StatusInternalServerError, err.Error())
	// 	}
	// 	return c.JSON(http.StatusOK, trx)
	// } else {
	// 	//expense
	// 	var accountID int
	// 	if req.BudgetID != nil {
	// 		//check that account should be under projects
	// 		isValid := false
	// 		for _, v := range budgetAccountIDs {
	// 			if v == *req.BudgetID {
	// 				//good
	// 				isValid = true
	// 			}
	// 		}
	// 		if !isValid {
	// 			return c.JSON(http.StatusInternalServerError, "Invalid budget ID")
	// 		}
	// 		//harus transalate dari budgetID ke accountID
	// 		isValid = false
	// 		for _, budget := range proj.Budgets {
	// 			if budget.ID == uint(*req.BudgetID) {
	// 				accountID = int(budget.AccountID)
	// 				isValid = true
	// 			}
	// 		}
	// 		if !isValid {
	// 			return c.JSON(http.StatusInternalServerError, "Invalid budget account ID")
	// 		}
	// 	} else {
	// 		return c.JSON(http.StatusInternalServerError, "Expense account must be defined")
	// 		// accountID = int(projAccountID)
	// 	}
	// 	trx, err := h.transactionStore.CreateTransfer(int(projAccountID), accountID, uint(math.Abs(float64(req.Amount))), req.Remarks, req.TransactionDate, false)
	// 	// trx, err := h.transactionStore.CreateTransaction(accountID, req.Amount, req.Remarks, req.SoD, req.TransactionDate)
	// 	if err != nil {
	// 		return c.JSON(http.StatusInternalServerError, err.Error())
	// 	}

	return c.JSON(http.StatusOK, trx)
	// }

}

// func (h *Handler) CreateProjectTransfer(c echo.Context) error {
// 	//todo: refactor this, bussiness logic should be nowhere here

// 	req := &createProjectTransferRequest{}
// 	var sourceAccount int
// 	var targetAccount int
// 	if err := req.bind(c); err != nil {
// 		return c.JSON(http.StatusUnprocessableEntity, err.Error())
// 	}
// 	trfDir, isSameProject := req.analyze()
// 	if !isSameProject && trfDir != ProjectToProject {
// 		//illegal
// 		return c.JSON(http.StatusUnauthorized, "Pocket can only transfered to parent project")
// 	}
// 	if isSameProject && trfDir == ProjectToPocket {
// 		//illegal
// 		return c.JSON(http.StatusUnauthorized, "Cannot transfer project to pocket")
// 	}
// 	if isSameProject && trfDir == PocketToProject {
// 		//illegal
// 		return c.JSON(http.StatusUnauthorized, "Cannot transfer pocket to project")
// 	}

// 	//get project and its accounts
// 	_, projectSourceAccountID, _, err := h.projectStore.GetProjectDetails(req.ProjectIDSource)
// 	if err != nil {
// 		return c.JSON(http.StatusUnauthorized, err.Error())
// 	}
// 	_, projectTargetAccountID, _, err := h.projectStore.GetProjectDetails(req.ProjectIDTarget)
// 	if err != nil {
// 		return c.JSON(http.StatusUnauthorized, err.Error())
// 	}

// 	//sof: project
// 	if trfDir == ProjectToPocket || trfDir == ProjectToProject {
// 		_, _, _, err := h.projectStore.GetProjectDetails(req.ProjectIDSource)
// 		if err != nil {
// 			return c.JSON(http.StatusUnauthorized, "Invalid Source Project ID")
// 		}
// 		sourceAccount = int(projectSourceAccountID)
// 	}

// 	//sof: pocket
// 	if trfDir == PocketToProject || trfDir == PocketToPocket {
// 		//check budgetid validity
// 		bgt, err := h.projectStore.CheckBudgetIDValidity(int(*req.BudgetIDSource), req.ProjectIDSource)
// 		if err != nil {
// 			return c.JSON(http.StatusUnauthorized, "Invalid Source Budget ID")
// 		}
// 		sourceAccount = int(bgt.AccountID)
// 	}

// 	//tof: project
// 	if trfDir == PocketToProject || trfDir == ProjectToProject {
// 		_, _, _, err := h.projectStore.GetProjectDetails(req.ProjectIDTarget)
// 		if err != nil {
// 			return c.JSON(http.StatusUnauthorized, "Invalid Target Project ID")
// 		}
// 		targetAccount = int(projectTargetAccountID)
// 	}

// 	//tof: pocket
// 	if trfDir == ProjectToPocket || trfDir == PocketToPocket {
// 		bgt, err := h.projectStore.CheckBudgetIDValidity(int(*req.BudgetIDTarget), req.ProjectIDTarget)
// 		if err != nil {
// 			return c.JSON(http.StatusUnauthorized, "Invalid Target Budget ID")
// 		}

// 		targetAccount = int(bgt.AccountID)
// 	}

// 	ret, err := h.transactionStore.CreateTransfer(sourceAccount, targetAccount, req.Amount, "TRF: "+req.Remarks, req.TrxDate, true)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, err.Error())
// 	}
// 	return c.JSON(http.StatusOK, ret)
// }

func (h *Handler) UpdateProject(c echo.Context) error {
	req := &updateProjectRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	isOpen := false
	if req.Status == "open" {
		isOpen = true
	}
	edit := model.Project{
		BaseModel: model.BaseModel{
			ID: uint(req.ProjectID),
		},
		IsOpen:      isOpen,
		Description: req.Description,
	}
	err := h.projectStore.UpdateProject(edit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	prj, _, _, err := h.projectStore.GetProjectDetails(int(req.ProjectID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, prj)
}

func (h *Handler) CreateBankTransfer(c echo.Context) error {
	req := &createBankTransferRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	ret, err := h.transactionStore.CreateBankTransaction(req.SourceBankAccountID, req.DestBankAccountID, req.Amount, req.Remarks, req.Meta, req.Notes, req.TrxDate)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	return c.JSON(http.StatusOK, ret)
}

func (h *Handler) GetProjectAnalysis(c echo.Context) error {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	//rewrite commonrequest inject using context to pass id
	res, err := h.projectStore.GetProjectAnalysis(projectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if len(res) == 0 {
		return c.JSON(http.StatusNotFound, nil)

	}
	return c.JSON(http.StatusOK, res)
}
