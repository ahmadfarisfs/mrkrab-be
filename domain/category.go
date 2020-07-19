package domain

import (
	"context"
	//"github.com/jinzhu/gorm"
)

// Transaction ...
type Category struct {
	BaseModel
	Name string `gorm:"not null" json:"name" validate:"required"`
	Icon string `gorm:"null" json:"icon"`
}

// CategoryUsecase represent the Category's usecases (business process)
type CategoryUsecase interface {
	Fetch(ctx context.Context, limitPerPage int64, page int64, filter map[string]string) ([]Category, error)
	GetByID(ctx context.Context, id int64) (Category, error)
	Add(ctx context.Context, Category *Category) error
	Update(ctx context.Context, Category *Category) error
	Delete(ctx context.Context, id int64) error
}

// CategoryRepository represent the Categorys's repository contract -> implemented in db conn
type CategoryRepository interface {
	Fetch(ctx context.Context, limitPerPage int64, page int64) (res []Category, err error)
	GetByID(ctx context.Context, id int64) (Category, error)
	Update(ctx context.Context, ar *Category) error
	Store(ctx context.Context, a *Category) error
	Delete(ctx context.Context, id int64) error
}
