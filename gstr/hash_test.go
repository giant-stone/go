package gstr_test

import (
	"reflect"
	"testing"

	"github.com/giant-stone/go/gstr"
)

func TestMd5(t *testing.T) {
	for _, item := range []struct {
		s    string
		want string
	}{
		{"你好，世界", "dbefd3ada018615b35588a01e216ae6e"},
		{"hello", "5d41402abc4b2a76b9719d911017c592"},
		{"", "d41d8cd98f00b204e9800998ecf8427e"},
	} {
		got := gstr.Md5(item.s)
		if !reflect.DeepEqual(item.want, got) {
			t.Errorf("Md5 -%s- want %s got %s", item.s, item.want, got)
		}
	}
}

func TestMd5FromBytes(t *testing.T) {
	for _, item := range []struct {
		s    string
		want string
	}{
		{"你好，世界", "dbefd3ada018615b35588a01e216ae6e"},
		{"hello", "5d41402abc4b2a76b9719d911017c592"},
		{"", "d41d8cd98f00b204e9800998ecf8427e"},
	} {
		got := gstr.Md5FromBytes([]byte(item.s))
		if !reflect.DeepEqual(item.want, got) {
			t.Errorf("Md5 -%s- want %s got %s", item.s, item.want, got)
		}
	}
}

func TestSha1(t *testing.T) {
	for _, item := range []struct {
		s    string
		want string
	}{
		{"你好，世界", "3becb03b015ed48050611c8d7afe4b88f70d5a20"},
		{"hello world", "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"},
		{"", "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
	} {
		got := gstr.Sha1(item.s)
		if !reflect.DeepEqual(item.want, got) {
			t.Errorf("Sha1 -%s- want %s got %s", item.s, item.want, got)
		}
	}
}

func TestSha1FromBytes(t *testing.T) {
	for _, item := range []struct {
		s    string
		want string
	}{
		{"你好，世界", "3becb03b015ed48050611c8d7afe4b88f70d5a20"},
		{"hello world", "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"},
		{"", "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
	} {
		got := gstr.Sha1FromBytes([]byte(item.s))
		if !reflect.DeepEqual(item.want, got) {
			t.Errorf("Sha1 -%s- want %s got %s", item.s, item.want, got)
		}
	}
}
