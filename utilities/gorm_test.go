package utilities

import (
	"context"
	"testing"
	"time"

	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestGormCount(t *testing.T) {
	db, err := gorm.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.AutoMigrate(&domain.User{})

	err = db.Create(&domain.User{
		FirstName: "Bambcang",
		LastName:  "Sudibyo",
		//	Role:      "jendral",
	}).Error
	if err != nil {
		panic("Cannot create " + err.Error())
	}
	err = db.Create(&domain.User{
		FirstName: "Bambcang",
		LastName:  "Sudibyo",
		//		Role:      "jendral",
	}).Error
	if err != nil {
		panic("Cannot create " + err.Error())
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(2500)*time.Millisecond)
	defer cancel()
	val := 0
	res, errDB := CountWithContext(ctx, db.Model(&domain.User{}), &val)
	if errDB != nil {
		panic(" errorDB")
	}
	if res.Error != nil {
		panic("apfcj")
	}
	t.Log(val)

}
