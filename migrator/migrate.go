package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/ahmadfarisfs/mrkrab-be/domain"

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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Migrator().DropTable(&domain.User{})
	db.Migrator().DropTable(&domain.Project{})
	db.Migrator().DropTable(&domain.ProjectBudget{})
	db.Migrator().DropTable(&domain.Category{})
	db.Migrator().DropTable(&domain.Transaction{})

	//db.Migrator().DropTable("user_projects")
	err = db.Set("gorm:table_options", "ENGINE=InnoDB").
		AutoMigrate(&domain.Project{},
			&domain.User{}, &domain.ProjectBudget{}, &domain.Category{},
			&domain.Transaction{})
	//db.Migrator().CreateConstraint()
	if err != nil {
		panic(err)
	}
}
