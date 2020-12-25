package main

import (
	"github.com/ahmadfarisfs/krab-core/db"
	_ "github.com/ahmadfarisfs/krab-core/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/ahmadfarisfs/krab-core/handler"
	"github.com/ahmadfarisfs/krab-core/router"
	"github.com/ahmadfarisfs/krab-core/store"
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
		panic("Error Reading config")
	}
	r := router.New()

	//r.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := r.Group("/api")

	d := db.New()
	db.AutoMigrate(d)
	as := store.NewAccountStore(d)
	ts := store.NewTransactionStore(d)
	ps := store.NewProjectStore(d)
	us := store.NewUserStore(d)
	ms := store.NewMutationStore(ps, d)
	h := handler.NewHandler(as, ts, ps, us, ms)
	h.Register(v1)
	r.Logger.Fatal(r.Start("127.0.0.1:" + viper.GetString(`service.port`)))
}
