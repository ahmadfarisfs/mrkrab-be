package handler

import (
	"github.com/ahmadfarisfs/krab-core/utils"
	"github.com/labstack/echo/v4"
)

//Register will
func (h *Handler) Register(v1 *echo.Group) {
	//	v1.Use()
	//	v1.Use(utils.ParseCommonMiddleware)

	fAccountGroup := v1.Group("/account/financial")
	fAccountGroup.POST("", h.RegisterFinancialAccount)
	fAccountGroup.GET("/:id", h.ViewFinancialAccountSummary)
	fAccountGroup.GET("", h.ListFinancialAccount, utils.ParseCommonMiddleware)
	fAccountGroup.PUT("", h.UpdateFinancialAccount)

	bAccountGroup := v1.Group("/account/bank")
	bAccountGroup.POST("", h.RegisterBankAccount)
	bAccountGroup.GET("/:id", h.ViewBankAccountSummary)
	bAccountGroup.GET("", h.ListBankAccount, utils.ParseCommonMiddleware)
	bAccountGroup.PUT("", h.UpdateBankAccount)

	trxGroup := v1.Group("/transactions")
	trxGroup.POST("/income", h.CreateIncomeTransaction)
	trxGroup.POST("/expense", h.CreateExpenseTransaction)
	trxGroup.POST("/transfer/bank", h.CreateBankTransferTransaction)
	trxGroup.POST("/transfer/project", h.CreateProjectTransferTransaction)

	// trxGroup.PUT("/income", h.upda)
	// trxGroup.PUT("/expense", h.CreateExpenseTransaction)
	// trxGroup.PUT("/transfer/bank", h.CreateBankTransferTransaction)
	// trxGroup.PUT("/transfer/project", h.CreateProjectTransferTransaction)

	trxGroup.GET("/:id", h.ViewTransactionDetails)
	trxGroup.GET("", h.ListTransaction, utils.ParseCommonMiddleware)

	prjGroup := v1.Group("/projects")
	prjGroup.GET("/:id", h.GetProject)
	prjGroup.GET("/financial/:id", h.GetProjectAnalysis)

	prjGroup.GET("", h.ListProject)
	prjGroup.POST("", h.CreateProject)
	prjGroup.DELETE("/:id", h.DeleteProject)
	prjGroup.PUT("", h.UpdateProject)

	userGroup := v1.Group("/users")
	userGroup.GET("/:id", h.GetUser)
	userGroup.DELETE("/:id", h.DeleteUser)
	userGroup.GET("", h.ListUser)
	userGroup.POST("", h.CreateUser)

	payRecFroup := v1.Group("/payrec")
	payRecFroup.GET("", h.ListPayRec)
	payRecFroup.POST("", h.CreatePayRec)
	payRecFroup.PATCH("/approve/:id", h.Approve)
	payRecFroup.PATCH("/reject/:id", h.Reject)

	authGroup := v1.Group("/auth")
	authGroup.POST("/login", h.Login)
	authGroup.POST("/authenticate", h.Authenticate)

}
