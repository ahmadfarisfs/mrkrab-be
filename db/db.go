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
		&model.FinancialAccount{},
		&model.Transaction{},
		&model.BankAccountMutation{},
		&model.FinancialAccountMutation{},
		&model.Project{},
		&model.BankAccount{},
		&model.User{},
		&model.PayRec{},
	)

	if err != nil {
		panic("Error migration" + err.Error())
	}
	// res := model.FinancialAccount{}
	// err = db.Model(&model.FinancialAccount{}).Where("account_name = 'ACCOUNT-REVENUE'").First(&res).Error
	// if errors.Is(err, gorm.ErrRecordNotFound) {

	// 	err = db.Exec("INSERT INTO `accounts`(`id`,`account_name`) VALUES(0,'ACCOUNT-REVENUE')").Error
	// 	if err != nil {
	// 		panic("Error create revenue account" + err.Error())
	// 	}

	// 	err = db.Exec("UPDATE accounts SET ID=0 WHERE account_name='ACCOUNT-REVENUE'").Error
	// 	if err != nil {
	// 		panic("Error create revenue account update id" + err.Error())
	// 	}
	// } else {
	// 	//found, check id has to be 0
	// 	if res.ID != 0 {
	// 		err = db.Exec("UPDATE accounts SET ID=0 WHERE account_name='ACCOUNT-REVENUE'").Error
	// 		if err != nil {
	// 			panic("Error [2] create revenue account update id" + err.Error())
	// 		}
	// 	}
	// }

}
