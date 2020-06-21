package domain

import (
	"context"

	"gorm.io/gorm"
)

// User ...
type User struct {
	gorm.Model
	FirstName string    `gorm:"not null" json:"firstname" validate:"required"`
	LastName  string    `gorm:"not null" json:"lastname" validate:"required"`
	Email     string    `gorm:"not null;unique" json:"email" validate:"required,email"`
	Phone     string    `gorm:"not null;unique" json:"phone" validate:"required"`
	Photo     string    `json:"photo"`
	Role      string    `gorm:"not null;type:enum('sa','pic','member')" json:"role"`
	Password  string    `gorm:"not null" json:"-"`
	Projects  []Project `json:"projects" gorm:"many2many:user_projects;foreignkey:id;references:id;"`
}

// Belongs To: `ForeignKey` specifies foreign key field owned by the current model, `References` specifies the association's primary key
// Has One/Many: `ForeignKey` specifies foreign key for the association, `References` specifies the current model's primary key
// Many2Many: `ForeignKey` specifies the current model's primary key, `JoinForeignKey` specifies join table's foreign key that refers to `ForeignKey`
//            `References` specifies the association's primary key, `JoinReferences` specifies join table's foreign key that refers to `References`
// For multiple foreign keys, it can be separated by commas

// UserUsecase represent the user's usecases (business process)
type UserUsecase interface {
	Fetch(ctx context.Context, limitPerPage int64, page int64) ([]User, error)
	GetByID(ctx context.Context, id int64) (User, error)
	GetByRole(ctx context.Context, role string) ([]User, error)
	Update(ctx context.Context, user *User) error
	Register(context.Context, *User) error
	GetProjectInvolved(ctx context.Context, userID int64) (map[ProjectMemberRole][]Project, error)
	Delete(ctx context.Context, id int64) error
}

// UserRepository represent the users's repository contract -> implemented in db conn
type UserRepository interface {
	Fetch(ctx context.Context, limitPerPage int64, page int64) (res []User, err error)
	GetByID(ctx context.Context, id int64) (User, error)
	GetByRole(ctx context.Context, role string) ([]User, error)
	Update(ctx context.Context, ar *User) error
	Store(ctx context.Context, a *User) error
	Delete(ctx context.Context, id int64) error
}
