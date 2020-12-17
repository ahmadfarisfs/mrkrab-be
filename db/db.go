package db

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"

	"os"

	"github.com/ahmadfarisfs/krab-core/model"
	"gorm.io/gorm"
)

func New() *gorm.DB {
	dsn := "user=postgres password=root dbname=bankcore port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
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
	)
	if err != nil {
		panic(err)
	}
}
