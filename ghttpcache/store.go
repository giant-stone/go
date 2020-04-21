package ghttpcache

import (
	"errors"
)

var (
	ErrObjectNotExist = errors.New("store object not exist")
)

type Store interface {
	Create(key string, value []byte, expireAt int64) (err error)
	Read(key string, expireAt int64) (value []byte, err error)
	Delete(key string) (err error)
}
