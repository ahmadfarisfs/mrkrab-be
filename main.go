package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/ahmadfarisfs/krab-core/db"
	"github.com/ahmadfarisfs/krab-core/handler"
	"github.com/ahmadfarisfs/krab-core/router"
	"github.com/ahmadfarisfs/krab-core/store"
	"github.com/labstack/echo/v4"

	// "github.com/labstack/echo"
	"github.com/spf13/viper"
	// echo-swagger middleware
)

// @title Swagger Example API
// @version 1.0
// @description Conduit API
// @title Conduit API

// @host 127.0.0.1:8585
// @BasePath /api

// @schemes http https
// @produce	application/json
// @consumes application/json

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		if os.Getenv("MRKRAB_CONFIG") == "" {
			panic("Error Reading config from env")
		} else {
			//load from env if file not found
			err := viper.ReadConfig(strings.NewReader(os.Getenv("MRKRAB_CONFIG")))
			if err != nil {
				panic("error read config from ENV2: " + err.Error())
			}

		}
	}
	r := router.New()
	r.GET("/echo", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Hello from mrkrabs")
	})

	r.GET("/money/echo", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Hello from mrkrabs money")
	})
	//r.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := r.Group("/api")

	d := db.New()
	db.AutoMigrate(d)
	fa := store.NewFinancialAccountStore(d)
	as := store.NewBankAccountStore(d)
	ts := store.NewTransactionStore(d)
	ps := store.NewProjectStore(d)
	us := store.NewUserStore(d)
	ms := store.NewMutationStore(ps, d)
	// prs := store.NewPayRecStore(ts, ps, d)
	h := handler.NewHandler(fa, as, ts, ps, us, ms)
	h.Register(v1)
	r.Logger.Fatal(r.Start("0.0.0.0:" + viper.GetString(`service.port`)))
}
