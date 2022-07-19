package gstr_test

import (
	"testing"

	"github.com/giant-stone/go/gstr"
	"github.com/stretchr/testify/require"
)

func TestRemoveHtmlTag(t *testing.T) {
	samples := []struct {
		s    string
		want string
	}{
		{
			"十五&amp;小宇 打架（恩爱）日常 老父亲吃瓜（狗粮）<em class=\"keyword\">看戏</em>ing",
			"十五&amp;小宇 打架（恩爱）日常 老父亲吃瓜（狗粮）看戏ing",
		},

		{
			"生化危机6-猎杀特工4月19日，<em class=\"keyword\">摸鱼看戏</em>",
			"生化危机6-猎杀特工4月19日，摸鱼看戏",
		},

		{
			"<h1>hello </h1>world",
			"hello world",
		},

		{
			"<em class=\"keyword\">电影最TOP</em>：他，定义了演技",
			"电影最TOP：他，定义了演技",
		},
	}

	for _, item := range samples {
		got := gstr.RemoveHtmlTag(item.s)
		require.Equal(t, item.want, got, item.s)
	}
}
