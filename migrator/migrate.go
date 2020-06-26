package main

import (
	"bytes"
	"fmt"
	"log"
	"net/url"

	"github.com/ahmadfarisfs/mrkrab-be/domain"
	"github.com/ahmadfarisfs/mrkrab-be/utilities"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {

	localConfigFilename := "config.json"
	if utilities.FileExists(localConfigFilename) {
		log.Println("Found Local Config!")
		viper.SetConfigFile(localConfigFilename)
	} else {
		log.Println("Not Found Local Config, finding on google secret manager!")
		secrets := utilities.FetchSecret("silmioti", "projects/633186564272/secrets/MySQL-Config/latest")
		viper.SetConfigType("json")
		viper.ReadConfig(bytes.NewBuffer(secrets))
	}

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	dbConn.Migrator().DropTable(&domain.User{})
	dbConn.Migrator().DropTable(&domain.Project{})
	dbConn.Migrator().DropTable(&domain.ProjectBudget{})
	dbConn.Migrator().DropTable(&domain.Category{})
	dbConn.Migrator().DropTable(&domain.Transaction{})
	dbConn.Migrator().DropTable("user_projects")
	dbConn.Migrator().DropTable("project_members")

	err = dbConn.Set("gorm:table_options", "ENGINE=InnoDB").
		AutoMigrate(&domain.Project{}, &domain.User{}, &domain.Category{},
			&domain.ProjectBudget{},
			&domain.Transaction{})
	if err != nil {
		panic(err)
	}
	dbConn.Create(&domain.User{
		Email:     "aku@komo.com",
		FirstName: "aku",
		LastName:  "komo",
		Phone:     "0855",
		Role:      "sa",
	})

	picid := 1
	dbConn.Create(&domain.Project{
		Name:        "proyeku",
		ProjectType: "onetime",
		Status:      "ongoing",
		PICID:       &picid,
	})
	dbConn.Create(&domain.Category{
		Name: "Makanan",
	})
	dbConn.Create(&domain.ProjectBudget{
		ProjectID:  1,
		CategoryID: 1,
		Amount:     10000,
	})
	if err != nil {
		panic(err)
	}
	log.Println("Success!")
}
