package ghttpcache

import (
	"context"
	"log"
	"time"

	"github.com/giant-stone/go/ghttp"
)

// HttpCache is a HTTP caching on local implementation.
type HttpCache struct {
	// Debug enable logging in verbose
	Debug bool

	ExpireDuration time.Duration

	HttpTimeout time.Duration
	HttpProxy   string

	Coder Coder

	Store Store
}

func New(c Coder, s Store) (rs *HttpCache) {
	return &HttpCache{
		Debug:          false,
		ExpireDuration: time.Duration(24) * time.Hour,
		HttpTimeout:    time.Duration(10) * time.Second,
		Coder:          c,
		Store:          s,
	}
}

func (its *HttpCache) Set(key string, value []byte, expireAt int64) (err error) {
	valueEncoded, err := its.Coder.Compress(value)
	if err != nil {
		return
	}

	return its.Store.Create(key, valueEncoded, expireAt)
}

func (its *HttpCache) Get(key string, expireAt int64) (rs []byte, err error) {
	value, err := its.Store.Read(key, expireAt)
	if err != nil {
		return
	}

	return its.Coder.Decompress(value)
}

func (its *HttpCache) Fetch(fullUrl string) (rs []byte, err error) {
	rq := ghttp.New(context.Background()).
		SetTimeout(its.HttpTimeout).
		SetRequestMethod("GET").
		SetUri(fullUrl).
		SetRandomUserAgent(true).
		SetProxy(its.HttpProxy)

	err = rq.Send()
	if ghttp.CheckRequestErr(fullUrl, rq.RespStatus, rq.RespBody, err) {
		return
	}

	if len(rq.RespBody) > 0 {
		rs = rq.RespBody
	}
	return
}

func (its *HttpCache) GetOrFetch(key string, expireAt int64) (rs []byte, err error) {
	v, err := its.Get(key, expireAt)
	if err != nil {
		if err == ErrObjectNotExist {
			log.Printf("[debug] cache miss key -%s-", key)

			rs, err = its.Fetch(key)
			if err != nil {
				return
			}

			errSet := its.Set(key, rs, expireAt)
			if errSet != nil {
				log.Printf("[warning] Set key -%s- expireAt=%d %v", key, expireAt, errSet)
			}
		}
	} else {
		log.Printf("[debug] cache hit key -%s-", key)
		rs = v
	}

	return
}

func (its *HttpCache) Hit(key string) (rs bool) {
	if _, err := its.Store.Read(key, 0); err == nil {
		rs = true
	}
	return
}

func (its *HttpCache) Delete(key string) (err error) {
	return its.Store.Delete(key)
}
