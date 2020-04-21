package ghttpcache_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/giant-stone/go/gutil"

	"github.com/giant-stone/go/ghttpcache"
)

const (
	sampleUrl = "https://tools.ietf.org/rfc/rfc1945.txt"
)

func ExampleHttpCache() {
	var err error

	c := ghttpcache.New(ghttpcache.NewCoderGzip(), ghttpcache.NewStoreMem())
	c.Debug = true
	c.ExpireDuration = time.Duration(50) * time.Millisecond

	key := sampleUrl
	value, err := c.GetOrFetch(key, 0)
	gutil.CheckErr(err)
	log.Printf("[debug] len(body)=%d", len(value))

	rs := c.Hit(key)
	log.Printf("[debug] get key -%s- hit=%t", key, rs)

	err = c.Delete(key)
	gutil.CheckErr(err)
	log.Printf("[debug] delete key -%s-", key)

	rs = c.Hit(key)
	log.Printf("[debug] get key -%s- hit=%t", key, rs)
}

func TestHttpCache_Get(t *testing.T) {
	timeout := time.Duration(50) * time.Millisecond

	c := ghttpcache.New(ghttpcache.NewCoderGzip(), ghttpcache.NewStoreMem())
	expireAt := time.Now().Add(timeout).UTC().Unix()
	key := sampleUrl
	value, err := c.GetOrFetch(key, expireAt)
	if err != nil {
		t.Errorf("want err=nil got %v", err)
	}
	if len(value) == 0 {
		t.Errorf("want len(value)>0 got 0")
	}

	time.Sleep(timeout)
	time.Sleep(timeout)

	rs := c.Hit(key)
	if rs != true {
		t.Errorf("want hit=false got true")
	}
}

func TestHttpCache_Set(t *testing.T) {
	timeout := time.Duration(50) * time.Millisecond

	c := ghttpcache.New(ghttpcache.NewCoderGzip(), ghttpcache.NewStoreMem())
	expireAt := time.Now().Add(timeout).UTC().Unix()
	key := sampleUrl

	rs := c.Hit(key)
	if rs != false {
		t.Errorf("want hit=false got true")
	}

	var err error
	value := []byte(`hello world`)

	err = c.Set(key, value, expireAt)
	if err != nil {
		t.Errorf("want err=nil got %v", err)
	}
	if len(value) == 0 {
		t.Errorf("want len(value)>0 got 0")
	}

	valueGot, _ := c.Get(key, expireAt)
	strWant := string(value)
	strGot := string(valueGot)
	if strWant != strGot {
		t.Errorf("want str=-%s- got -%s-", strWant, strGot)
	}

	rs = c.Hit(key)
	if rs != true {
		t.Errorf("want hit=true got false")
	}
}

func TestMain(m *testing.M) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	os.Exit(m.Run())
}
