package server

import (
	"bytes"
	"crypto/aes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/rahul0tripathi/clipd/services/keybind"
	"github.com/rahul0tripathi/clipd/services/keychain"
	"github.com/rahul0tripathi/clipd/services/pasteboard"
	"github.com/rahul0tripathi/clipd/util"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/pbkdf2"
	"regexp"
	"runtime"
	"sync"
)

var (
	DomainSeparator = "!"
	// DecryptRegex this should compile
	DecryptRegex = regexp.MustCompile(`^(?P<desc>\[.*])!(?P<permutationKey>[a-zA-Z0-9-]+)!(?P<encryptedText>[a-zA-Z0-9]+)$`)
)

type Server interface {
	Listen() error
	Shutdown() error
}

type localServer struct {
	mu                 *sync.Mutex
	errChan            chan error
	saltHash           string
	keyPath            string
	logger             *zap.Logger
	keychainService    keychain.SecretsService
	keyBind            keybind.Service
	pasteboardService  pasteboard.PasteboardService
	permutationService util.PermutationService
}

func (l *localServer) init() (err error) {
	// opening db in read mode only
	db, dbErr := util.NewLevelDB(l.keyPath, true, l.logger)
	if dbErr != nil {
		return dbErr
	}
	l.permutationService, err = util.NewDefaultPermutationService(db, l.logger)
	if err != nil {
		return
	}
	l.keychainService, err = keychain.NewKeychainService(l.logger)
	if err != nil {
		return
	}
	l.keyBind, err = keybind.NewDefaultKeybindService(l.logger)
	if err != nil {
		return
	}
	l.pasteboardService, err = pasteboard.NewDefaultPasteboardService()
	return err
}
func (l *localServer) aesWrapper(input []byte, permutedKey *[32]byte, encrypt bool) (output []byte, err error) {
	masterKey, getPassErr := l.keychainService.GetPassword(nil)
	if getPassErr != nil {
		l.logger.Error("failed to get password stored in  keychain", zap.Error(err))
		return nil, getPassErr
	}
	l.logger.Debug("aesWrapper:", zap.String("master-key", string(masterKey)))
	finalKeyPbkdf2 := pbkdf2.Key(masterKey, permutedKey[:], 10000, 32, sha256.New)
	l.logger.Debug("aesWrapper:", zap.String("final-key", fmt.Sprintf("%x", finalKeyPbkdf2)))
	finalKey := [32]byte{}
	copy(finalKey[:], finalKeyPbkdf2[:])
	if encrypt {
		output, err = util.Encrypt(input, &finalKey)
	} else {
		output, err = util.Decrypt(input, &finalKey)
	}
	return
}
func (l *localServer) formatEncrypted(key string, data string) string {
	return fmt.Sprintf("[WHOAMI]%s%s%s%x", DomainSeparator, key, DomainSeparator, data)
}
func (l *localServer) decrypt() {
	// idc if it blows up
	if isLocked := l.mu.TryLock(); !isLocked {
		l.logger.Error("failed to acquire lock in encrypt")
		return
	}
	defer func() {
		l.mu.Unlock()
		runtime.GC()
	}()
	clipdString, err := l.pasteboardService.ReadFromPasteboard()
	match := DecryptRegex.FindStringSubmatch(string(clipdString))
	if len(match) != 4 {
		l.logger.Error("encryptedText does not match format")
		return
	}
	key := match[2]
	if key == "" {
		l.logger.Error("permutation key not found")
		return
	}
	encryptedText := match[3]
	if encryptedText == "" {
		l.logger.Error("encryptedText key not found")
		return
	}
	if len(encryptedText) < aes.BlockSize {
		l.logger.Error("encryptedText length less than min block size")
		return
	}
	l.logger.Debug("decrypt:", zap.String("plainText", encryptedText), zap.Int("length", len(encryptedText)))
	if err != nil {
		l.logger.Error("failed to read from pasteboard", zap.Error(err))
		return
	}
	encryptedTextParsedHex, err := hex.DecodeString(encryptedText)
	if err != nil {
		l.logger.Error("failed to DecodeString hex", zap.Error(err))
		return
	}
	salt := [32]byte{}
	copy(salt[:], l.saltHash)
	permutedSalt, err := l.permutationService.PermuteWithKey(key, salt)
	if err != nil {
		l.logger.Error("failed to compute permuted salt", zap.Error(err))
		return
	}
	l.logger.Debug("decrypt:", zap.String("permuted-key", fmt.Sprintf("%x", permutedSalt[:])))
	decoded, err := l.aesWrapper(encryptedTextParsedHex, &permutedSalt, false)
	if err != nil {
		l.logger.Error("failed to encrypt", zap.Error(err))
		return
	}
	l.logger.Debug("decrypt:", zap.String("decrypted", string(decoded)))
	l.logger.Debug("decrypt:", zap.String("formatted-string", string(decoded)))
	err = l.pasteboardService.WriteToPasteboard(decoded)
	if err != nil {
		l.logger.Error("failed to WriteToPasteboard", zap.Error(err))
		return
	}
}
func (l *localServer) encrypt() {
	// shamelessly using the TryLock idc if it blows up
	if isLocked := l.mu.TryLock(); !isLocked {
		l.logger.Error("failed to acquire lock in encrypt")
		return
	}
	defer func() {
		l.mu.Unlock()
		runtime.GC()
	}()
	plainText, err := l.pasteboardService.ReadFromPasteboard()
	if len(plainText) < aes.BlockSize {
		plainText = append(plainText, bytes.Repeat([]byte("~"), aes.BlockSize-len(plainText))...)
	}
	l.logger.Debug("encrypt:", zap.String("plainText", string(plainText)), zap.Int("length", len(plainText)))
	if err != nil {
		l.logger.Error("failed to read from pasteboard", zap.Error(err))
		return
	}
	l.logger.Debug("encrypt:", zap.String("plainText-base64", string(plainText)), zap.Int("length", len(plainText)))
	permutedSalt := [32]byte{}
	copy(permutedSalt[:], l.saltHash)
	permutedCtx, err := l.permutationService.PermuteRandomly(permutedSalt)
	if err != nil {
		l.logger.Error("failed to compute permuted salt", zap.Error(err))
		return
	}
	l.logger.Debug("encrypt:", zap.String("permutation-key", permutedCtx.Key), zap.String("permuted-key", fmt.Sprintf("%x", permutedCtx.Data[:])))
	encrypted, err := l.aesWrapper(plainText, &permutedCtx.Data, true)
	if err != nil {
		l.logger.Error("failed to encrypt", zap.Error(err))
		return
	}
	formattedString := l.formatEncrypted(permutedCtx.Key, string(encrypted))
	l.logger.Debug("encrypt:", zap.String("formatted-string", formattedString))
	err = l.pasteboardService.WriteToPasteboard([]byte(formattedString))
	if err != nil {
		l.logger.Error("failed to WriteToPasteboard", zap.Error(err))
		return
	}
}
func (l *localServer) Listen() error {
	doneChan := make(chan struct{})
	decryptListener, err := l.keyBind.GetDecryptListener()
	if err != nil {
		return err
	}
	encryptListener, err := l.keyBind.GetEncryptListener()
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-decryptListener:
				l.decrypt()
			case err := <-l.errChan:
				l.logger.Error("decryptListener caught an error, closing", zap.Error(err))
				select {
				case doneChan <- struct{}{}:
				default:
				}
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case <-encryptListener:
				l.encrypt()
			case err := <-l.errChan:
				l.logger.Error("encryptListener caught an error, closing", zap.Error(err))
				select {
				case doneChan <- struct{}{}:
				default:
				}
				return
			}
		}
	}()
	<-doneChan
	return nil
}

func (l *localServer) Shutdown() error {
	select {
	case l.errChan <- errors.New("graceful shutdown called"):
	default:
		return errors.New("failed to shutdown server")
	}
	return nil
}
func NewLocalServer(saltHash string, keyPath string, logger *zap.Logger) (Server, error) {
	server := &localServer{
		logger:   logger,
		saltHash: saltHash,
		keyPath:  keyPath,
		errChan:  make(chan error),
		mu:       &sync.Mutex{},
	}
	err := server.init()
	if err != nil {
		return nil, err
	}
	return server, nil
}
