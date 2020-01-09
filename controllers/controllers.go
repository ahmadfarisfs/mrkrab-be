package controllers

import (
	"mrkrab-be/db"
	"net/http"
	"obb-new-parking/controllers"

	"github.com/thedevsaddam/renderer"
)

var DB = db.ConnectDatabase()

var Rnd = renderer.New()

func GenerateStandardResponse(w http.ResponseWriter, reason string, statusCode int) {
	controllers.Rnd.JSON(w, statusCode, struct {
		Reason string `json:"reason"`
	}{
		Reason: reason,
	})

}
