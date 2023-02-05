package util

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"go.uber.org/zap"
)

type StorageEngine interface {
	Put(key string, data []byte) error
	Get(key string) ([]byte, error)
}

type LocalLevelDB struct {
	path       string
	dbInstance *leveldb.DB
	logger     *zap.Logger
}

func (l LocalLevelDB) Put(key string, data []byte) error {
	return l.dbInstance.Put([]byte(key), data, &opt.WriteOptions{
		Sync: true,
	})
}

func (l LocalLevelDB) Get(key string) ([]byte, error) {
	return l.dbInstance.Get([]byte(key), nil)
}

func NewLevelDB(path string, readOnly bool, logger *zap.Logger) (StorageEngine, error) {
	var err error
	ldb := &LocalLevelDB{
		path:   path,
		logger: logger,
	}
	logger.Info("initialising localStorageEngine", zap.String("path", path))
	ldb.dbInstance, err = leveldb.OpenFile(path, &opt.Options{
		ReadOnly: readOnly,
	})
	if err != nil {
		return nil, err
	}
	return ldb, nil
}
