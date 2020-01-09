package controllers

import (
	"mrkrab-be/controllers"
	model "mrkrab-be/models"
	utils "mrkrab-be/utils"
	"net/http"
)

func NewUser(w http.ResponseWriter, r *http.Request) {
	req := model.User{}
	err := utils.DecodeJSONBody(w, r, &req)
	if err != nil {
		controllers.GenerateStandardResponse(w, err.Error(), http.StatusInternalServerError)
	}
	err = controllers.DB.Create(&req).Error
	if err != nil {
		controllers.GenerateStandardResponse(w, err.Error(), http.StatusInternalServerError)
	}
	controllers.GenerateStandardResponse(w, "success", http.StatusCreated)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	return
}

func RemoveUser(w http.ResponseWriter, r *http.Request) {
	return
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	return
}
