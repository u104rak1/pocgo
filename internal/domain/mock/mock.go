//go:build !production

package mock

// The command to create a mock of an entity file is as follows:
// mockgen -source=internal/domain/{entity}/{entity}_repository.go -destination=internal/domain/mock/mock_{entity}_repository.go -package=mock

// The command to create a mock of a service file is as follows:
// mockgen -source=internal/domain/{service}/{service}_service.go -destination=internal/domain/mock/mock_{service}_service.go -package=mock
