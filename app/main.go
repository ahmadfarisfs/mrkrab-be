package main

import (
	"fmt"
	"log"
	"net/url"
	"time"

	userHandler "github.com/ahmadfarisfs/mrkrab-be/user/delivery/http"
	userRepo "github.com/ahmadfarisfs/mrkrab-be/user/repository/mysql"
	userUsecase "github.com/ahmadfarisfs/mrkrab-be/user/usecase"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	viper.SetConfigFile(`config.json`)
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
	log.Println("here")
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	e := echo.New()
	//	middL := _articleHttpDeliveryMiddleware.InitMiddleware()
	//	e.Use(middL.CORS)

	//repo init
	userRP := userRepo.NewUserRepo(dbConn)

	//usecase init
	userUC := userUsecase.NewUserUsecase(userRP, nil, timeoutContext)

	//handler init
	userHandler.NewUserHandler(e, userUC)

	log.Fatal(e.Start(viper.GetString("server.address")))
}
