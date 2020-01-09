package main

import (
	"log"
	"mrkrab-be/controllers"
	model "mrkrab-be/models"
	"mrkrab-be/routes"
	"net/http"
	"os"
	"time"

	//	"github.com/ddo/go-mux-mvc/models/logger"
	"github.com/joho/godotenv"
	// init db
)

const (
	defaultPort = "8008"

	idleTimeout       = 30 * time.Second
	writeTimeout      = 180 * time.Second
	readHeaderTimeout = 10 * time.Second
	readTimeout       = 10 * time.Second
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	err := controllers.DB.AutoMigrate(&model.Transaction{}, &model.User{}, &model.TransactionCategory{}, &model.Project{}).Error
	if err != nil {
		panic("Cannot migrate DB: " + err.Error())
	}
	server := &http.Server{
		Addr:              "0.0.0.0:" + port,
		IdleTimeout:       idleTimeout,
		WriteTimeout:      writeTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
	}
	log.Println("Serving on port " + port)
	http.Handle("/", routes.Handlers())
	err = server.ListenAndServe()
	if err != nil {
		log.Println("ERR ListenAndServe:", err)
	}
}
