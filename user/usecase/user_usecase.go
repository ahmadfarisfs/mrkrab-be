package usecase

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ahmadfarisfs/mrkrab-be/domain"
	"golang.org/x/crypto/bcrypt"
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
func (u *userUseCase) Login(ctx context.Context, userName string, password string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	user, err := u.userRepo.GetByUsername(ctx, userName)
	if err != nil {
		return domain.User{}, err
	}
	if comparePasswords(user.Password, password) {
		return user, nil
	}
	return domain.User{}, errors.New("Password Mismatch")
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
	user.Password = hashAndSalt(user.Password)
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

func hashAndSalt(pwd string) string {

	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}
func comparePasswords(hashedPwd string, plainPwd string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plainPwd))
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
