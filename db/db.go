package db

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"

	"os"

	"github.com/ahmadfarisfs/krab-core/model"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func New() *gorm.DB {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	//dsn := "user=postgres password=root dbname=bankcore port=5432 sslmode=disable TimeZone=Asia/Jakarta" &loc=Asia%2FJakarta
	//db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
	dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Asia%2FJakarta"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second, // Slow SQL threshold
				LogLevel:      logger.Info, // Log level
				Colorful:      true,        // Disable color
			},
		),
	})

	if err != nil {
		fmt.Println("storage err: ", err)
	}

	return db
}

//AutoMigrate migrate all model except trigger
func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.Account{},

		&model.Transaction{},
		&model.Mutation{},
		&model.Project{},
		&model.Budget{},
		&model.User{},
	)
	if err != nil {
		panic("Error migration" + err.Error())
	}
}
