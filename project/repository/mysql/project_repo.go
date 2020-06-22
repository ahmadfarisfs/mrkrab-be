package mysql

import (
	"context"
	"log"

	"github.com/ahmadfarisfs/mrkrab-be/domain"
	"github.com/ahmadfarisfs/mrkrab-be/utilities"

	"gorm.io/gorm"
)

type mysqlProjectRepo struct {
	DB *gorm.DB
}

func NewProjectRepo(db *gorm.DB) domain.ProjectRepository {
	return &mysqlProjectRepo{
		DB: db,
	}
}
func (r *mysqlProjectRepo) GetProjectsByUser(ctx context.Context, userID int64) (map[domain.ProjectMemberRole][]domain.Project, error) {
	panic("not imepl")
}
func (r *mysqlProjectRepo) RemoveMember(ctx context.Context, projectID int64, userID int64) error {
	panic("not imepl")
}

func (r *mysqlProjectRepo) GetProjectMember(ctx context.Context, projectID int64) (map[domain.ProjectMemberRole][]domain.User, error) {
	panic("not imepl")
}

func (r *mysqlProjectRepo) Fetch(ctx context.Context, limitPerPage int64, page int64) (res []domain.Project, totalRecord int, totalPage int, err error) {
	users := []domain.Project{}
	//query := r.DB.Find(&users)

	pagingInfo := utilities.Paging(ctx, &utilities.Param{
		DB:      r.DB.Model(&domain.User{}), //.Where("id > ?", 0),
		Page:    int(page),
		Limit:   int(limitPerPage),
		OrderBy: []string{"id desc"},
	}, &users)
	log.Println(users)
	return users, pagingInfo.TotalRecord, pagingInfo.TotalPage, err
}

func (r *mysqlProjectRepo) GetByID(ctx context.Context, id int64) (domain.Project, error) {
	user := domain.Project{}
	err := r.DB.Where("id = ?", id).First(&user).Error
	return user, err
}

func (r *mysqlProjectRepo) GetByRole(ctx context.Context, role string) ([]domain.Project, error) {
	user := []domain.Project{}
	err := r.DB.Where("role = ?", role).Find(&user).Error
	return user, err
}
func (r *mysqlProjectRepo) AddMember(ctx context.Context, projectID int64, userID int64, role domain.ProjectMemberRole) error {
	panic("Not implemented")
	//return r.DB.Save(ar).Error
}
func (r *mysqlProjectRepo) Update(ctx context.Context, ar *domain.Project) error {
	return r.DB.Save(ar).Error
}
func (r *mysqlProjectRepo) Store(ctx context.Context, a *domain.Project) error {
	return r.DB.Create(a).Error
}
func (r *mysqlProjectRepo) Delete(ctx context.Context, id int64) error {
	return r.DB.Where("id = ?", id).Delete(&domain.Project{}).Error
}
