//go:build !production

package mock

// The command to create a mock of a usecase file is as follows:
// mockgen -source=internal/application/{domain}/{fileName}.go -destination=internal/application/mock/mock_{fileName}.go -package=mock
