package authentication

type IJWTService interface {
	GenerateAccessToken(userID string) (string, error)
	GetUserIDFromAccessToken(accessToken string) (string, error)
}
