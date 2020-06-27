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
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(viper.GetString(`jwt.secret`)),
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.RequestURI()
			log.Println(path)
			if strings.Contains(path, "login") {
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

	//usecase init
	userUC := userUsecase.NewUserUsecase(userRP, nil, timeoutContext)
	projectUC := projectUsecase.NewProjectUseCase(projectRP, timeoutContext)

	//handler init
	userHandler.NewUserHandler(e, userUC)
	projectHandler.NewProjectHandler(e, projectUC)

	log.Fatal(e.Start(viper.GetString("server.address")))
}
