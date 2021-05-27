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
	ret := model.User{Fullname: name, Username: username, Role: role, Email: email, Password: string(hashedPassword)}
	err = ps.db.Model(&model.User{}).Create(&ret).Error
	if err != nil {
		return model.User{}, err
	}

	//create account
	err = ps.db.Model(&model.FinancialAccount{}).Create(&model.FinancialAccount{
		AccountName: "USER-" + strings.ToUpper(username) + "-" + strconv.Itoa(int(time.Now().Unix())),
	}).Error
	if err != nil {
		return model.User{}, err
	}

	return ret, err
}
func (ps *UserStore) ListUser(req utils.CommonRequest) ([]model.User, int, error) {
	ret := []model.User{}
	var count int64
	//query builder
	initQuery := ps.db

	err := initQuery.Model(&model.User{}).Count(&count).Error
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

func (ps *UserStore) GetUserDetails(id int) (model.User, error) {
	usr := model.User{}
	err := ps.db.Model(&model.User{}).Where("id= ?", id).First(&usr).Error
	return usr, err
}
