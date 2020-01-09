package model

import "database/sql/driver"

type UserRoleType int64

const (
	SystemAdmininstrator UserRoleType = iota
	Treasurer
	Engineer
	BoardOfDirector
)

var types = [...]string{
	"ADM",
	"TRS",
	"ENG",
	"BOD",
}

func (t UserRoleType) String() string {
	return types[t]
}

func (u *UserRoleType) Scan(value interface{}) error {
	*u = UserRoleType(value.(int64))
	return nil
}

func (u UserRoleType) Value() (driver.Value, error) { return int64(u), nil }
