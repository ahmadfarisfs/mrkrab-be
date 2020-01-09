package main

import (
	"log"
	"mrkrab-be/db"
	"mrkrab-be/migrator"
	model "mrkrab-be/models"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	db := db.ConnectDatabase()
	//db, err := gorm.Open("sqlite3", "db.db")
	//	if err != nil {
	//		panic("fail to connect dbase")
	//	}
	defer db.Close()
	err := db.AutoMigrate(&model.User{}).Error
	if err != nil {
		panic(err)
	}
	newUser := model.User{
		Email:             "myemail@email.com",
		Role:              model.Engineer,
		Name:              "Engineer 01",
		Password:          "Mypassword",
		Phone:             80989999,
		ProfilePictureURL: "urls",
		Username:          time.Now().Format(time.RFC1123),
	}
	err = db.Create(&newUser).Error
	if err != nil {
		log.Println("PANIC CREATE")
		panic(err)
	}
	user := model.User{}

	err = db.Where("name = ?", "Engineer 01").First(&user).Error

	if err != nil {
		log.Println("PANIC WHERE")

		panic(err)
	}
	log.Println(user)

	migrator.Migrate(db, &model.Project{})
	time.Sleep(time.Second * 1000)
	//	db.AutoMigrate(&model.Transaction{}, &model.User{}, &model.TransactionCategory{}, &model.Project{})
}
