package usecase

import (
	"context"
	"time"

	"github.com/ahmadfarisfs/mrkrab-be/domain"
)

type userUseCase struct {
	userRepo    domain.UserRepository
	projectRepo domain.ProjectRepository
	timeout     time.Duration
}

func NewUserUsecase(a domain.UserRepository, p domain.ProjectRepository, timeout time.Duration) domain.UserUsecase {
	return &userUseCase{
		projectRepo: p,
		userRepo:    a,
		timeout:     timeout,
	}
}

func (u *userUseCase) Fetch(ctx context.Context, limitPerPage int64, page int64) (users []domain.User, totalRecord int, totalPage int, err error) {
	if limitPerPage == 0 {
		limitPerPage = 10
	}
	if page == 0 {
		page = 1
	}
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	return u.userRepo.Fetch(ctx, limitPerPage, page)
}

func (u *userUseCase) GetByID(ctx context.Context, id int64) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	return u.userRepo.GetByID(ctx, id)
}

func (u *userUseCase) GetByRole(ctx context.Context, role string) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	return u.userRepo.GetByRole(ctx, role)
}

func (u *userUseCase) Update(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	return u.userRepo.Update(ctx, user)
}

func (u *userUseCase) Register(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	return u.userRepo.Store(ctx, user)
}

func (u *userUseCase) GetProjectInvolved(ctx context.Context, userID int64) (map[domain.ProjectMemberRole][]domain.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	return u.projectRepo.GetProjectsByUser(ctx, userID)
}

func (u *userUseCase) Delete(ctx context.Context, id int64) error {
	return nil
}
