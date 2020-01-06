package main

import (
	model "mrkrab-be/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	db, err := gorm.Open("sqlite3", "db.db")
	if err != nil {
		panic("fail to connect dbase")
	}
	defer db.Close()
	db.AutoMigrate(&model.Transaction{}, &model.User{}, &model.TransactionCategory{}, &model.Project{})
}
