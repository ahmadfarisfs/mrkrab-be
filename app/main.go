package main

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"time"

	//"github.com/ahmadfarisfs/mrkrab-be/middleware"
	//secretmanager "cloud.google.com/go/secretmanager/apiv1"
	//secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	projectHandler "github.com/ahmadfarisfs/mrkrab-be/project/delivery/http"
	userHandler "github.com/ahmadfarisfs/mrkrab-be/user/delivery/http"
	"github.com/ahmadfarisfs/mrkrab-be/utilities"

	projectRepo "github.com/ahmadfarisfs/mrkrab-be/project/repository/mysql"
	projectUsecase "github.com/ahmadfarisfs/mrkrab-be/project/usecase"
	userRepo "github.com/ahmadfarisfs/mrkrab-be/user/repository/mysql"
	userUsecase "github.com/ahmadfarisfs/mrkrab-be/user/usecase"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {

	localConfigFilename := "config.json"
	if utilities.FileExists(localConfigFilename) {
		viper.SetConfigFile(localConfigFilename)
	} else {
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
		log.Fatal(err)
	}
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	e := echo.New()
	e.Use(middleware.CORS())

	//repo init
	userRP := userRepo.NewUserRepo(dbConn)
	projectRP := projectRepo.NewProjectRepo(dbConn)

	//usecase init
	userUC := userUsecase.NewUserUsecase(userRP, nil, timeoutContext)
	projectUC := projectUsecase.NewProjectUseCase(projectRP, timeoutContext)

	//handler init
	userHandler.NewUserHandler(e, userUC)
	projectHandler.NewProjectHandler(e, projectUC)

	log.Fatal(e.Start(viper.GetString("server.address")))
}
