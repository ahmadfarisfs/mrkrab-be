package db

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

func ConnectDatabase() *gorm.DB {
	//Load environmenatal variables
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mySQLAddress := os.Getenv("MYSQL_SERVER")
	mySQLPass := os.Getenv("MYSQL_PASSWORD")
	mySQLUser := os.Getenv("MYSQL_USERNAME")
	databaseName := os.Getenv("MYSQL_DATABASE")
	db, err := gorm.Open("mysql", mySQLUser+":"+mySQLPass+"@("+mySQLAddress+")/"+databaseName+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	return db
}
