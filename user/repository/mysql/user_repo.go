package mysql

import (
	"context"

	"github.com/ahmadfarisfs/mrkrab-be/domain"

	"gorm.io/gorm"
	//	"github.com/jinzhu/gorm"
	//"github.com/jinzhu/gorm"
)

type mysqlUserRepo struct {
	DB *gorm.DB
}

func (r *mysqlUserRepo) Fetch(ctx context.Context, limitPerPage int64, page int64) (res []domain.User, err error) {
	users := []domain.User{}
	query := r.DB
	/*	utilities.Paging(ctx, &utilities.Param{
			DB:      query,
			Page:    int(page),
			Limit:   int(limitPerPage),
			OrderBy: []string{"id desc"},
		}, &users)
	*/return users, nil
}
func (r *mysqlUserRepo) GetByID(ctx context.Context, id int64) (domain.User, error) {
	user := domain.User{}
	err := r.DB.Where("id = ?", id).First(&user).Error
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
