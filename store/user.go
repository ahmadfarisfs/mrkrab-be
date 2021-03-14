package store

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ahmadfarisfs/krab-core/model"
	"github.com/ahmadfarisfs/krab-core/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (ps *UserStore) CreateUser(name string, username string, password string, email string, role string) (model.User, error) {
	//TODO: add username check, email check, revive the dead account

	// Salt and hash the password using the bcrypt algorithm
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return model.User{}, err
	}
	//TODO: use transaction
	tx := ps.db.Begin()
	//create account
	account := model.Account{
		AccountType: "BANK", //create bank account for user
		AccountName: "USER-" + strings.ToUpper(username) + "-" + strconv.Itoa(int(time.Now().Unix())),
	}
	err = tx.Model(&model.Account{}).Create(&account).Error
	if err != nil {
		tx.Rollback()
		return model.User{}, err
	}
	ret := model.User{
		AccountID: account.ID,
		Fullname:  name, Username: username, Role: role, Email: email, Password: string(hashedPassword)}
	err = tx.Model(&model.User{}).Create(&ret).Error
	if err != nil {
		tx.Rollback()
		return model.User{}, err
	}

	return ret, tx.Commit().Error
}

func (ps *UserStore) ListUser(req utils.CommonRequest) ([]model.User, int, error) {
	ret := []model.User{}
	var count int64
	//query builder
	initQuery := ps.db

	err := initQuery.Model(&model.User{}).Preload("Account").Count(&count).Error
	if err != nil {
		return ret, int(count), err
	}
	log.Println(req)

	//actually fetch data with limit and offset
	quer := utils.AppendCommonRequest(initQuery, req)
	err = quer.Find(&ret).Error
	return ret, int(count), err
}

func (ps *UserStore) DeleteUser(id int) error {
	return ps.db.Where("id = ?", id).Delete(&model.User{}).Error
}
