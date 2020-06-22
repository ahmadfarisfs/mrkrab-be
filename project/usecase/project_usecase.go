package usecase

import (
	"context"
	"time"

	"github.com/ahmadfarisfs/mrkrab-be/domain"
	"gorm.io/gorm"
)

type mysqlProjectUseCase struct {
	DB          *gorm.DB
	projectRepo domain.ProjectRepository
	timeout     time.Duration
}

func NewProjectUseCase(db *gorm.DB, p domain.ProjectRepository, timeout time.Duration) domain.ProjectRepository {
	return &mysqlProjectUseCase{
		projectRepo: p,
		DB:          db,
		timeout:     timeout,
	}
}

func (p *mysqlProjectUseCase) Fetch(ctx context.Context, limitPerPage int64, page int64) (res []domain.Project, totalRecord int, totalPage int, err error) {
	if limitPerPage == 0 {
		limitPerPage = 10
	}
	if page == 0 {
		page = 1
	}
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()
	return p.projectRepo.Fetch(ctx, limitPerPage, page)

}
func (p *mysqlProjectUseCase) GetByID(ctx context.Context, id int64) (domain.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()
	return p.projectRepo.GetByID(ctx, id)
}
func (p *mysqlProjectUseCase) Update(ctx context.Context, ar *domain.Project) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()
	return p.projectRepo.Update(ctx, ar)
}
func (p *mysqlProjectUseCase) Store(ctx context.Context, a *domain.Project) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()
	return p.projectRepo.Store(ctx, a)
}
func (p *mysqlProjectUseCase) Delete(ctx context.Context, id int64) error {
	return nil
}
func (p *mysqlProjectUseCase) IsPICAssigned(ctx context.Context, projectID int64) (bool, error) {
	panic("Note implemented")
}
func (p *mysqlProjectUseCase) GetProjectsByUser(ctx context.Context, userID int64) (map[domain.ProjectMemberRole][]domain.Project, error) {
	panic("Note implemented")
}
func (p *mysqlProjectUseCase) GetProjectMember(ctx context.Context, projectID int64) (map[domain.ProjectMemberRole][]domain.User, error) {
	panic("Note implemented")
}
func (p *mysqlProjectUseCase) AddMember(ctx context.Context, projectID int64, userID int64, role domain.ProjectMemberRole) error {
	panic("Note implemented")
}
func (p *mysqlProjectUseCase) RemoveMember(ctx context.Context, projectID int64, userID int64) error {
	panic("Note implemented")
}
