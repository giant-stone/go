package ghttpcache

import (
	"log"
	"time"
)

type StoreObject interface {
	GetData() (rs []byte)
	GetExpireAt() (rs int64)
	SetExpireAt(e int64) (err error)
	HasExpired() (rs bool)
}

type storeMem struct {
	data map[string]StoreObject
}

type storeObjectMem struct {
	data     []byte
	expireAt int64
}

func NewStoreObjectMem(data []byte, expireAt int64) *storeObjectMem {
	return &storeObjectMem{data: data, expireAt: expireAt}
}

func (its *storeObjectMem) GetData() (rs []byte) {
	return its.data
}

func (its *storeObjectMem) HasExpired() (rs bool) {
	return time.Now().UTC().Unix() < its.expireAt
}

func (its *storeObjectMem) GetExpireAt() (rs int64) {
	return its.expireAt
}

func (its *storeObjectMem) SetExpireAt(e int64) (err error) {
	its.expireAt = e
	return
}

func (its *storeMem) Read(key string, expireAt int64) (value []byte, err error) {
	v, found := its.data[key]
	if !found {
		err = ErrObjectNotExist
		return
	}

	if v.GetExpireAt() == 0 && expireAt > 0 {
		v.SetExpireAt(expireAt)
	}

	if v.HasExpired() {
		errDelete := its.Delete(key)
		if errDelete != nil {
			log.Printf("[warning] Delete key -%s- %v", key, errDelete)
		}
		err = ErrObjectNotExist
		return
	}

	value = v.GetData()
	return
}

func (its *storeMem) Delete(key string) (err error) {
	delete(its.data, key)
	return nil
}

func (its *storeMem) Create(key string, value []byte, expireAt int64) (err error) {
	its.data[key] = NewStoreObjectMem(value, expireAt)
	return nil
}

func NewStoreMem() *storeMem {
	m := map[string]StoreObject{}
	return &storeMem{
		data: m,
	}
}
