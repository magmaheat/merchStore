package hasher

//go:generate go run github.com/vektra/mockery/v2@latest --name=PasswordHasher --output=./../../internal/mocks
type PasswordHasher interface {
	Hash(password string) (string, error)
	Check(password, hash string) bool
}
