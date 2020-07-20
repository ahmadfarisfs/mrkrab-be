package domain

import (
	"context"
	"encoding/json"
	"strings"
	"time"
)

// first create a type alias
type JSONDateOnly time.Time

// User ...
type User struct {
	//gorm.Model
	BaseModel
	Username  string  `gorm:"not null;unique" json:"username" validate:"required"`
	FirstName string  `gorm:"not null" json:"firstname" validate:"required"`
	LastName  string  `gorm:"not null" json:"lastname" validate:"required"`
	Birthday  string  `gorm:"not null;column:birthday;type:date" json:"birthday" validate:"required"`
	Email     string  `gorm:"not null;unique" json:"email" validate:"required,email"`
	Phone     string  `gorm:"not null;unique;type:varchar(255)" json:"phone" validate:"required"`
	Photo     *string `json:"photo"`
	Role      string  `gorm:"not null;type:enum('sa','pic','member','secretary')" json:"role" validate:"required"`
	Password  string  `gorm:"not null" json:"password" validate:"required"`
	//	Projects  []Project `json:"projects" gorm:"many2many:user_projects;foreignkey:id;references:id;"`
}

// imeplement Marshaler und Unmarshalere interface
func (j *JSONDateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = JSONDateOnly(t)
	return nil
}

func (j JSONDateOnly) MarshalJSON() ([]byte, error) {
	return json.Marshal(j)
}

// Maybe a Format function for printing your date
func (j JSONDateOnly) Format(s string) string {
	t := time.Time(j)
	return t.Format(s)
}

// Belongs To: `ForeignKey` specifies foreign key field owned by the current model, `References` specifies the association's primary key
// Has One/Many: `ForeignKey` specifies foreign key for the association, `References` specifies the current model's primary key
// Many2Many: `ForeignKey` specifies the current model's primary key, `JoinForeignKey` specifies join table's foreign key that refers to `ForeignKey`
//            `References` specifies the association's primary key, `JoinReferences` specifies join table's foreign key that refers to `References`
// For multiple foreign keys, it can be separated by commas

// UserUsecase represent the user's usecases (business process)
type UserUsecase interface {
	Login(ctx context.Context, userName string, password string) (User, error)
	Fetch(ctx context.Context, limitPerPage int64, page int64) (users []User, totalRecord int, totalPage int, err error)
	GetByID(ctx context.Context, id int64) (User, error)
	GetByRole(ctx context.Context, role string) ([]User, error)
	Update(ctx context.Context, user *User) error
	Register(context.Context, *User) error
	GetProjectInvolved(ctx context.Context, userID int64) (map[ProjectMemberRole][]Project, error)
	Delete(ctx context.Context, id int64) error
}

// UserRepository represent the users's repository contract -> implemented in db conn
type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (User, error)
	Fetch(ctx context.Context, limitPerPage int64, page int64) (res []User, totalRecord int, totalPage int, err error)
	GetByID(ctx context.Context, id int64) (User, error)
	GetByIDs(ctx context.Context, id []int64) ([]User, error)
	GetByRole(ctx context.Context, role string) ([]User, error)
	Update(ctx context.Context, ar *User) error
	Store(ctx context.Context, a *User) error
	Delete(ctx context.Context, id int64) error
}
