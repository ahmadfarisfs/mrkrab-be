package contract

type AuthStore interface {
	Login(username, password string) error
	// CreateAccount(name string, parentID *uint) (model.Account, error)
	// GetAccountDetails(id int) (model.Account, error)
}
