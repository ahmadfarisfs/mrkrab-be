package usecase

import (
	"context"
	"time"

	"github.com/ahmadfarisfs/mrkrab-be/domain"
	//	"gorm.io/gorm"
)

type projectUseCase struct {
	//	DB          *gorm.DB
	projectRepo domain.ProjectRepository
	timeout     time.Duration
}

func NewProjectUseCase(p domain.ProjectRepository, timeout time.Duration) domain.ProjectUsecase {
	return &projectUseCase{
		projectRepo: p,
		//	DB:          db,
		timeout: timeout,
	}
}

func (p *projectUseCase) Fetch(ctx context.Context, limitPerPage int64, page int64, filter map[string]string) (res []domain.Project, totalRecord int, totalPage int, err error) {
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
func (p *projectUseCase) GetByID(ctx context.Context, id int64) (domain.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()
	return p.projectRepo.GetByID(ctx, id)
}
func (p *projectUseCase) Update(ctx context.Context, ar *domain.Project) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()
	return p.projectRepo.Update(ctx, ar)
}
func (p *projectUseCase) Add(ctx context.Context, a *domain.Project) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()
	return p.projectRepo.Store(ctx, a)
}
func (p *projectUseCase) Delete(ctx context.Context, id int64) error {
	return nil
}
func (p *projectUseCase) AssignPIC(ctx context.Context, projectID int64, userID int64) error {
	panic("Note implemented")
}
func (p *projectUseCase) GetProjectsByUser(ctx context.Context, userID int64) (map[domain.ProjectMemberRole][]domain.Project, error) {
	panic("Note implemented")
}
func (p *projectUseCase) GetProjectMember(ctx context.Context, projectID int64) (map[domain.ProjectMemberRole][]domain.User, error) {
	panic("Note implemented")
}
func (p *projectUseCase) AssignMember(ctx context.Context, projectID int64, userID int64) error {
	panic("Note implemented")
}
func (p *projectUseCase) RemoveMember(ctx context.Context, projectID int64, userID int64) error {
	panic("Note implemented")
}

func (p *projectUseCase) SetStatus(ctx context.Context, projectID int64, status domain.ProjectStatus) error {
	panic("")
}
