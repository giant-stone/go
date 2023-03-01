package ghttp

import (
	"context"
	"net/http"
	"net/url"
	"testing"
)

func TestHttpRequest_GenerateRequest(t *testing.T) {
	tests := []struct {
		name string
		it   *HttpRequest
		want *http.Request
	}{
		{
			name: "succ",
			it:   &HttpRequest{Ctx: context.Background(), Uri: "https://ietf.org"},
			want: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Scheme: "https",
					Host:   "ietf.org",
					Path:   "/",
				},
			},
		},

		{
			name: "invalid URL thats contains space",
			it:   &HttpRequest{Uri: " https://ietf.org"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.it.GenerateRequest()
			if (tt.want == nil && tt.want != got) || (tt.want != nil && got == nil) {
				t.Errorf("HttpRequest.GenerateRequest() got %v, want %v", got, tt.want)
			}
		})
	}
}
