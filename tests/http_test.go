package tests

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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

func TestHttpUploadProxy(t *testing.T) {
	// 测试反向代理
	tests := []struct {
		name     string
		url      string
		method   string
		filepath string
		want     int
	}{
		{
			name:     "test proxy",
			url:      "http://localhost:7256/api/v1/mini-upload",
			method:   http.MethodPost,
			filepath: "../config.yaml",
			want:     http.StatusOK,
		},
		{
			name:     "test proxy",
			url:      "http://localhost:7256/api/v1/mini-upload",
			method:   http.MethodPost,
			filepath: "../mini-gateway.exe",
			want:     http.StatusRequestEntityTooLarge,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := &http.Client{}

			// 上传文件
			file := tt.filepath
			var bf bytes.Buffer
			writer := multipart.NewWriter(&bf)
			w, _ := writer.CreateFormFile("file", filepath.Base(tt.filepath))

			f, err := os.Open(file)
			if err != nil {
				t.Errorf("Open() error = %v", err)
				return
			}
			defer f.Close()
			n, err := io.Copy(w, f)
			if err != nil {
				t.Errorf("Copy() error = %v", err)
				return
			}
			log.Printf("copy %d bytes\n", n)
			writer.Close()

			req, _ := http.NewRequest(tt.method, tt.url, &bf)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			res, err := httpClient.Do(req)
			if err != nil {
				t.Errorf("Do() error = %v", err)
				return
			}

			if res.StatusCode == tt.want {
				t.Logf("upload success, status code: %d", res.StatusCode)
			} else {
				t.Errorf("upload failed, status code: %d", res.StatusCode)
			}
		})
	}
}
