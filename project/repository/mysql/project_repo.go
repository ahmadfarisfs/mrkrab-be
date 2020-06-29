package mysql

import (
	"context"

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
func (r *mysqlProjectRepo) RemoveMember(ctx context.Context, project domain.Project, user domain.User) error {
	err := r.DB.Model(&project).Association("Members").Delete(user)
	return err
}

func (r *mysqlProjectRepo) GetProjectMember(ctx context.Context, project domain.Project) ([]domain.User, error) {
	users := []domain.User{}
	err := r.DB.Model(&project).Association("Role").Find(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *mysqlProjectRepo) Fetch(ctx context.Context, limitPerPage int64, page int64) (res []domain.Project, totalRecord int, totalPage int, err error) {
	projects := []domain.Project{}

	pagingInfo := utilities.Paging(ctx, &utilities.Param{
		DB:      r.DB.Joins("PIC").Model(&domain.Project{}), //.Preload("Users"),
		Page:    int(page),
		Limit:   int(limitPerPage),
		OrderBy: []string{"id desc"},
	}, &projects)

	return projects, pagingInfo.TotalRecord, pagingInfo.TotalPage, err
}

func (r *mysqlProjectRepo) GetByID(ctx context.Context, id int64) (domain.Project, error) {
	project := domain.Project{}
	err := r.DB.Preload("Budgets").Preload("Members").Joins("PIC").Where("projects.id = ?", id).First(&project).Error
	return project, err
}

func (r *mysqlProjectRepo) GetByRole(ctx context.Context, role string) ([]domain.Project, error) {
	user := []domain.Project{}
	err := r.DB.Where("role = ?", role).Find(&user).Error
	return user, err
}
func (r *mysqlProjectRepo) AddMember(ctx context.Context, project domain.Project, users []domain.User) error {
	project.Members = users
	err := r.DB.Model(&project).Association("Members").Append(users)
	return err
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
