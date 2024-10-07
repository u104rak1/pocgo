package signup_usecase

import (
	account_domain "github.com/ucho456job/pocgo/internal/domain/account"
	user_domain "github.com/ucho456job/pocgo/internal/domain/user"
)

type SignupDTO struct {
	User        UserDTO
	AccessToken string
}

type UserDTO struct {
	ID       string
	Name     string
	Email    string
	Accounts []AccountDTO
}

type AccountDTO struct {
	ID            string
	Name          string
	Balance       float64
	Currency      string
	LastUpdatedAt string
}

func newSignupDTO(user *user_domain.User, account *account_domain.Account, accessToken string) *SignupDTO {
	return &SignupDTO{
		User: UserDTO{
			ID:    user.ID(),
			Name:  user.Name(),
			Email: user.Email(),
			Accounts: []AccountDTO{
				{
					ID:            account.ID(),
					Name:          account.Name(),
					Balance:       account.Balance().Amount(),
					Currency:      account.Balance().Currency(),
					LastUpdatedAt: account.LastUpdatedAt().String(),
				},
			},
		},
		AccessToken: accessToken,
	}
}
