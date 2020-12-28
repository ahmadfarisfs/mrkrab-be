package handler

import (
	"github.com/labstack/echo/v4"
)

//Register will
func (h *Handler) Register(v1 *echo.Group) {
	//	v1.Use()
	//	v1.Use(utils.ParseCommonMiddleware)
	accountGroup := v1.Group("/account")
	accountGroup.POST("", h.RegisterAccount)
	accountGroup.GET("/:id", h.ViewAccountSummary)
	accountGroup.GET("", h.ListAccount)
	//accountGroup.GET("/mutation", h.ViewMutation)

	trxGroup := v1.Group("/transactions")
	//shortcut - only valid for cash transaction (assets and expenses)
	trxGroup.POST("/create", h.CreateTransaction)
	trxGroup.GET("/:id", h.ViewTransactionDetails)
	trxGroup.GET("", h.ListMutation)

	trfGroup := v1.Group("/transfer")
	trfGroup.POST("", h.CreateTransfer)

	prjGroup := v1.Group("/projects")
	prjGroup.GET("/:id", h.GetProject)
	prjGroup.GET("", h.ListProject)
	prjGroup.POST("", h.CreateProject)
	prjGroup.DELETE("/:id", h.DeleteProject)
	prjGroup.PUT("", h.UpdateProject)

	prjGroup.POST("/pocket", h.CreatePocket)
	prjGroup.POST("/transaction", h.CreateProjectTransaction)
	prjGroup.POST("/transfer", h.CreateProjectTransfer)

	userGroup := v1.Group("/users")
	userGroup.GET("/:id", h.GetUser)
	userGroup.DELETE("/:id", h.DeleteUser)
	userGroup.GET("", h.ListUser)
	userGroup.POST("", h.CreateUser)

	payRecFroup := v1.Group("/payrec")
	// payRecFroup.GET("/:id", h.GetUser)
	// payRecFroup.DELETE("/:id", h.DeleteUser)
	payRecFroup.GET("", h.ListPayRec)
	payRecFroup.POST("", h.CreatePayRec)
	payRecFroup.GET("/approve/:id", h.Approve)
	payRecFroup.GET("/reject/:id", h.Reject)

}
