package tests

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestHttpProxy(t *testing.T) {
	// 测试反向代理
	tests := []struct {
		name   string
		url    string
		method string
		want   string
	}{
		// TODO: Add test cases.
		{
			name:   "test proxy",
			url:    "http://localhost:7256/api/v1/ping",
			method: http.MethodGet,
			want:   "pong",
		},
		{
			name:   "test proxy",
			url:    "http://localhost:7256/api/v1/user",
			method: http.MethodGet,
			want:   "GET",
		},
		{
			name:   "test proxy",
			url:    "http://localhost:7256/api/v1/user",
			method: http.MethodPost,
			want:   "POST",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := &http.Client{}
			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Errorf("NewRequest() error = %v", err)
				return
			}
			res, err := httpClient.Do(req)
			if err != nil {
				t.Errorf("Do() error = %v", err)
				return
			}
			defer res.Body.Close()
			b, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("ReadAll() error = %v", err)
				return
			}
			got := string(b)
			if !strings.Contains(got, tt.want) {
				t.Errorf("Get() body = %v, want %v", got, tt.want)
				return
			}
		})
	}
}
