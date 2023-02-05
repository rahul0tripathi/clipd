package keychain

type SecretsService interface {
	Init() error
	CreateNewSecret() error
	GetPassword() ([]byte, error)
}

type keychainService struct {
}

func NewKeychainService() (SecretsService, error) { return nil, nil }
