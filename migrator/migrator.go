package migrator

import (
	"log"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type ForeignKeyMigration struct {
	ForeignKey   string
	ForeignKeyOf string
	OnDelete     string
	OnUpdate     string
}

var Wg sync.WaitGroup

type TriggerFunc func() string

func AddTriggers(db *gorm.DB, triggers ...TriggerFunc) {
	if len(triggers) > 0 {
		for _, tr := range triggers {
			res := db.Exec(tr())
			if res.Error != nil {
				log.Println(res.Error)
			}
		}
	}
}

func Migrate(db *gorm.DB, model interface{}) {
	mod := reflect.Indirect(reflect.ValueOf(model))
	numfield := mod.NumField()
	fks := make([]ForeignKeyMigration, 0)

	for i := 0; i < numfield; i++ {
		fieldStruct := mod.Type().Field(i)
		fieldName := fieldStruct.Name
		fk := ForeignKeyMigration{}

		tag := fieldStruct.Tag.Get("migrator")

		if tag == "" {
			continue
		}

		log.Println("Tag", fieldName, "of", mod.Type().Name(), "has migrator: ", tag)

		alltags := strings.Split(tag, ";")
		for _, pair := range alltags {
			keyval := strings.Split(pair, ":")
			switch keyval[0] {
			case "foreignkey":
				fk.ForeignKey = keyval[1]
			case "foreignkeyof":
				fk.ForeignKeyOf = keyval[1]
			case "ondelete":
				fk.OnDelete = keyval[1]
			case "onupdate":
				fk.OnUpdate = keyval[1]
			}
		}

		fks = append(fks, fk)
	}

	Wg.Add(1)
	go func(fkey []ForeignKeyMigration) {
		if len(fkey) > 0 {
			for _, fk := range fkey {
				startOfArg := strings.Index(fk.ForeignKeyOf, "(")
				dependsOn := fk.ForeignKeyOf[:startOfArg]
				log.Println("foreign key for", fk.ForeignKey, "depends on", dependsOn)

				for {
					if db.HasTable(dependsOn) {
						break
					}
					log.Println(dependsOn, "not created.. waiting..")
					time.Sleep(100 * time.Millisecond)
					log.Println("checking again..")
				}

				log.Println(dependsOn, "created, Migrating..")
				res := db.AutoMigrate(model).AddForeignKey(fk.ForeignKey, fk.ForeignKeyOf, fk.OnDelete, fk.OnUpdate)
				if res.Error != nil {
					log.Println(res.Error)
				}
			}
		} else {
			db.AutoMigrate(model)
		}
		Wg.Done()
	}(fks)
}
