package mysql

import (
	"context"
	"log"

	"github.com/ahmadfarisfs/mrkrab-be/domain"
	"github.com/ahmadfarisfs/mrkrab-be/utilities"

	"gorm.io/gorm"
	//	"github.com/jinzhu/gorm"
	//"github.com/jinzhu/gorm"
)

type mysqlUserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) domain.UserRepository {
	return &mysqlUserRepo{
		DB: db,
	}
}
func (r *mysqlUserRepo) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	user := domain.User{}
	err := r.DB.Where("username = ?", username).
		First(&user).Error
	return user, err
}
func (r *mysqlUserRepo) Fetch(ctx context.Context, limitPerPage int64, page int64) (res []domain.User, totalRecord int, totalPage int, err error) {
	users := []domain.User{}
	pagingInfo := utilities.Paging(ctx, &utilities.Param{
		DB:      r.DB.Model(&domain.User{}),
		Limit:   int(limitPerPage),
		OrderBy: []string{"id asc"},
	}, &users)
	log.Println(users)
	return users, pagingInfo.TotalRecord, pagingInfo.TotalPage, err
}
func (r *mysqlUserRepo) GetByID(ctx context.Context, id int64) (domain.User, error) {
	user := domain.User{}
	err := r.DB.Where("id = ?", id).First(&user).Error
	return user, err
}

func (r *mysqlUserRepo) GetByIDs(ctx context.Context, id []int64) ([]domain.User, error) {
	user := []domain.User{}
	err := r.DB.Where("id IN (?)", id).Find(&user).Error
	return user, err
}

func (r *mysqlUserRepo) GetByRole(ctx context.Context, role string) ([]domain.User, error) {
	user := []domain.User{}
	err := r.DB.Where("role = ?", role).Find(&user).Error
	return user, err
}
func (r *mysqlUserRepo) Update(ctx context.Context, ar *domain.User) error {
	return r.DB.Save(ar).Error
}
func (r *mysqlUserRepo) Store(ctx context.Context, a *domain.User) error {
	return r.DB.Create(a).Error
}
func (r *mysqlUserRepo) Delete(ctx context.Context, id int64) error {
	return r.DB.Where("id = ?", id).Delete(&domain.User{}).Error
}
