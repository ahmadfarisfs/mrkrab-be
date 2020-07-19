package main

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	//"github.com/ahmadfarisfs/mrkrab-be/middleware"
	//secretmanager "cloud.google.com/go/secretmanager/apiv1"
	//secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"github.com/ahmadfarisfs/mrkrab-be/domain"
	projectHandler "github.com/ahmadfarisfs/mrkrab-be/project/delivery/http"
	transactionHandler "github.com/ahmadfarisfs/mrkrab-be/transaction/delivery/http"
	userHandler "github.com/ahmadfarisfs/mrkrab-be/user/delivery/http"

	"github.com/ahmadfarisfs/mrkrab-be/utilities"

	projectRepo "github.com/ahmadfarisfs/mrkrab-be/project/repository/mysql"
	projectUsecase "github.com/ahmadfarisfs/mrkrab-be/project/usecase"
	trxRepo "github.com/ahmadfarisfs/mrkrab-be/transaction/repository/mysql"
	trxUsecase "github.com/ahmadfarisfs/mrkrab-be/transaction/usecase"

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
	var err error
	if os.Getenv("GAE_APPLICATION") == "" {
		//run in local
		log.Println("Found Local File Config")
		viper.SetConfigFile("config.json")
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
			log.Printf("Defaulting to port %s", port)
		}
		err = viper.ReadInConfig()
	} else {
		//run in GAE
		log.Println("Not Found Local File Config, Searching on Google Secret Manager")
		secrets := utilities.FetchSecret("silmioti", "projects/silmioti/secrets/MySQL-Config/versions/latest")
		viper.SetConfigType("json")
		err = viper.ReadConfig(bytes.NewBuffer(secrets))
	}

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
	val.Add("loc", viper.GetString(`database.loc`))
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	err = dbConn.Set("gorm:table_options", "ENGINE=InnoDB").
		AutoMigrate(&domain.Project{}, &domain.User{}, &domain.Category{},
			&domain.ProjectBudget{},
			&domain.Transaction{})
	if err != nil {
		log.Fatal(err)
	}
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	e := echo.New()
	//	e.Use(middleware.Logger())
	//e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(viper.GetString(`jwt.secret`)),
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.RequestURI()
			hostname := c.Request().Host
			log.Println(hostname)
			if strings.Contains(path, "login") || strings.Contains(hostname, "localhost") {
				return true
			}
			return false
		},
	}))
	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, "Hello")
	})
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, "I am Healthy!")
	})

	//repo init
	userRP := userRepo.NewUserRepo(dbConn)
	projectRP := projectRepo.NewProjectRepo(dbConn)
	transactionRP := trxRepo.NewTransactionRepo(dbConn)

	//usecase init
	userUC := userUsecase.NewUserUsecase(userRP, projectRP, timeoutContext)
	projectUC := projectUsecase.NewProjectUseCase(projectRP, userRP, transactionRP, timeoutContext)
	transactionUD := trxUsecase.NewTransactionUseCase(transactionRP, userRP, timeoutContext)

	//handler init
	userHandler.NewUserHandler(e, userUC)
	projectHandler.NewProjectHandler(e, projectUC)
	transactionHandler.NewTransactionHandler(e, transactionUD)
	log.Fatal(e.Start(viper.GetString("server.address")))
}
