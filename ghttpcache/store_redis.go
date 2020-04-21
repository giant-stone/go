package ghttpcache

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type storeRedis struct {
	prefix string
	pool   *redis.Pool
}

const (
	TTL_REPLY_NOT_EXISTS   = -1
	TTL_REPLY_NEVER_EXPIRE = -2
)

func (its *storeRedis) Read(key string, expireAt int64) (value []byte, err error) {
	conn := its.pool.Get()
	defer conn.Close()

	value, err = redis.Bytes(conn.Do("GET", its.prefix+key))
	if err != nil {
		if err == redis.ErrNil {
			err = ErrObjectNotExist
			return
		}
	}

	if expireAt > 0 {
		ttl, errTtl := redis.Int(conn.Do("TTL", its.prefix+key))
		if errTtl != nil {
			err = errTtl
			return
		}

		if ttl == TTL_REPLY_NOT_EXISTS {
			err = ErrObjectNotExist
		} else if ttl == TTL_REPLY_NEVER_EXPIRE {
			_, err = conn.Do("EXPIREAT", its.prefix+key, expireAt)
		}
	}

	return
}

func (its *storeRedis) Delete(key string) (err error) {
	conn := its.pool.Get()
	defer conn.Close()

	_, err = conn.Do("DEL", its.prefix+key)
	return
}

func (its *storeRedis) Create(key string, value []byte, expireAt int64) (err error) {
	conn := its.pool.Get()
	defer conn.Close()

	if expireAt > 0 {
		ttl := expireAt - time.Now().UTC().Unix()
		_, err = conn.Do("SETEX", its.prefix+key, ttl, value)
	} else {
		_, err = conn.Do("SET", its.prefix+key, value)
	}
	return
}

// NewStoreRedis ...
//   url in IANA scheme (https://www.iana.org/assignments/uri-schemes/prov/redis)
func NewStoreRedis(url string) *storeRedis {
	poolRedis := redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.DialURL(url) },
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	return &storeRedis{
		prefix: "ghttpcache:",
		pool:   &poolRedis,
	}
}
