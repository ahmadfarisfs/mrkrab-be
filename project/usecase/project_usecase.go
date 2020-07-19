package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/ahmadfarisfs/mrkrab-be/domain"
)

type projectUseCase struct {
	projectRepo     domain.ProjectRepository
	userRepo        domain.UserRepository
	timeout         time.Duration
	transactionRepo domain.TransactionRepository
}

func NewProjectUseCase(p domain.ProjectRepository, u domain.UserRepository, t domain.TransactionRepository, timeout time.Duration) domain.ProjectUsecase {
	return &projectUseCase{
		transactionRepo: t,
		projectRepo:     p,
		userRepo:        u,
		timeout:         timeout,
	}
}

func (p *projectUseCase) FetchTransaction(ctx context.Context, limitPerPage int64, page int64, filter map[string]string) (res []domain.Transaction, totalRecord int, totalPage int, err error) {
	if limitPerPage == 0 {
		limitPerPage = 10
	}
	if page == 0 {
		page = 1
	}

	return nil, 0, 0, nil
}

func (p *projectUseCase) AddTransaction(ctx context.Context, projectID int64, trx domain.Transaction) error {
	if trx.ProjectID != int(projectID) {
		return errors.New("Project ID mismatch")
	}
	return p.transactionRepo.Store(ctx, &trx)
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
func (p *projectUseCase) GetProjectsByUser(ctx context.Context, userID int64) ([]domain.Project, error) {
	panic("Note implemented")
}
func (p *projectUseCase) GetProjectMember(ctx context.Context, projectID int64) ([]domain.User, error) {
	panic("Note implemented")

}
func (p *projectUseCase) SetBudget(ctx context.Context, projectID int64, cat domain.Category, amountLimit int) error {

	return nil
}

func (p *projectUseCase) AssignMember(ctx context.Context, projectID int64, userID []int64) error {
	//check project ID
	project, err := p.GetByID(ctx, projectID)
	if err != nil {
		return err
	}
	//check project status
	if project.Status == "close" {
		return errors.New("Cannot add member to closed project")
	}
	//check if member == pic
	for _, val := range userID {
		if project.PICID == val {
			return errors.New("One of assigned members is PIC in this project")
		}
	}

	//check user ids
	users, err := p.userRepo.GetByIDs(ctx, userID)
	if err != nil {
		return err
	}
	if len(users) != len(userID) {
		return errors.New("Cannot find one or more member")
	}
	//add member
	err = p.projectRepo.AddMember(ctx, project, users)
	if err != nil {
		return err
	}
	return nil
}
func (p *projectUseCase) RemoveMember(ctx context.Context, projectID int64, userID int64) error {
	//	panic("Not implemented")
	//check project ID
	project, err := p.GetByID(ctx, projectID)
	if err != nil {
		return err
	}
	//check user ids
	user, err := p.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	//add member
	err = p.projectRepo.RemoveMember(ctx, project, user)
	if err != nil {
		return err
	}
	return nil
}

func (p *projectUseCase) SetStatus(ctx context.Context, projectID int64, status domain.ProjectStatus) error {
	panic("")
}
