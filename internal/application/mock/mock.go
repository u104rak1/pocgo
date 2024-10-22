//go:build !production

package mock

// The command to create a mock of a usecase file is as follows:
// mockgen -source=internal/application/authentication/signup_usecase.go -destination=internal/application/mock/mock_signup_usecase.go -package=mock
