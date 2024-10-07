package signup_usecase

type SignupCmd struct {
	User    UserCmd
	Account AccountCmd
}

type UserCmd struct {
	Name     string
	Email    string
	Password string
}

type AccountCmd struct {
	Name     string
	Password string
	Balance  float64
	Currency string
}
