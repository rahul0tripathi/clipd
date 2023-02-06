package keychain

import (
	"errors"
	"github.com/keybase/go-keychain"
	"go.uber.org/zap"
	"os/user"
)

const (
	ChainService = "go-clipd"
	AccessGroup  = "clipd.service"
)

var (
	ErrSecretAlreadyExists = errors.New("secret already exists")
	ErrCreateSecret        = errors.New("failed to create a new secret")
	ErrGetSecret           = errors.New("failed to get secret")
	ErrSecretNotFound      = errors.New("secret not found in keychain")
)

type SecretsService interface {
	CreateNewSecret(password []byte, username []byte) error
	GetPassword(username []byte) ([]byte, error)
}

type keychainService struct {
	logger *zap.Logger
}

func (k *keychainService) CreateNewSecret(password []byte, username []byte) error {
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(ChainService)
	current, err := user.Current()
	if username == nil {
		if err != nil {
			return err
		}
		username = []byte(current.Username)
	}
	item.SetAccount(string(username))
	item.SetLabel("MasterPassword")
	item.SetAccessGroup(AccessGroup)
	item.SetData(password)
	item.SetSynchronizable(keychain.SynchronizableNo)
	item.SetAccessible(keychain.AccessibleWhenPasscodeSetThisDeviceOnly)
	err = keychain.AddItem(item)
	if err == keychain.ErrorDuplicateItem {
		return ErrSecretAlreadyExists
	}
	if err != nil {
		k.logger.Error("failed to create a new secret", zap.Error(err))
		return ErrCreateSecret
	}
	return nil
}
func (k *keychainService) GetPassword(username []byte) ([]byte, error) {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(ChainService)
	current, err := user.Current()
	if username == nil {
		if err != nil {
			return nil, err
		}
		username = []byte(current.Username)
	}
	query.SetAccount(string(username))
	query.SetAccessGroup(AccessGroup)
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)
	results, err := keychain.QueryItem(query)
	if err != nil {
		k.logger.Error("failed to query keychain items", zap.Error(err))
		return nil, ErrGetSecret
	} else {
		if len(results) == 0 {
			return nil, ErrSecretNotFound
		}
		return results[0].Data, nil
	}
}
func NewKeychainService(logger *zap.Logger) (SecretsService, error) {
	service := &keychainService{logger: logger}
	return service, nil
}
