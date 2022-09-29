package ghttp_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/giant-stone/go/ghttp"
	"github.com/giant-stone/go/glogging"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGhttpDoGhttpNew(t *testing.T) {
	glogging.Init(nil, "")

	ctrl := gomock.NewController(t)

	clientImpl := ghttp.NewMockHttpClient(ctrl)
	ghttp.UseImpl(clientImpl)

	wantBody := []byte(`hello world`)

	rq := ghttp.
		New().
		SetRequestMethod(ghttp.POST).
		SetUri("https://httpbin.org/post").
		SetPostBody(wantBody).GenerateRequest()

	clientImpl.EXPECT().Do(rq).Return(
		&http.Response{
			Body:       io.NopCloser(bytes.NewReader(wantBody)),
			StatusCode: http.StatusOK,
		}, nil,
	)

	rs, err := ghttp.New().Do(rq)
	require.NoError(t, err)

	gotBody, err := ghttp.ReadBody(rs)
	require.NoError(t, err)
	require.Equal(t, wantBody, gotBody)
}

func TestGhttpDoHttpNewRequest(t *testing.T) {
	glogging.Init(nil, "")

	ctrl := gomock.NewController(t)

	clientImpl := ghttp.NewMockHttpClient(ctrl)
	ghttp.UseImpl(clientImpl)

	wantBody := []byte(`hello world`)

	rq, _ := http.NewRequest("POST", "https://httpbin.org/post", io.NopCloser(bytes.NewReader(wantBody)))

	clientImpl.EXPECT().Do(rq).Return(
		&http.Response{
			Body:       io.NopCloser(bytes.NewReader(wantBody)),
			StatusCode: http.StatusOK,
		}, nil,
	)

	rs, err := ghttp.New().Do(rq)
	require.NoError(t, err)

	gotBody, err := ghttp.ReadBody(rs)
	require.NoError(t, err)
	require.Equal(t, wantBody, gotBody)
}
