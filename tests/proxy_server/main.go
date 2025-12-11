package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewResponse(code int, message string, data interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func (r *Response) Write(w http.ResponseWriter) {
	js, _ := json.Marshal(r)
	w.Write(js)
}

func main() {
	port := flag.Int("port", 5000, "server port")
	flag.Parse()

	log.Printf("server start at :%d", *port)
	http.HandleFunc("/api/v1/user", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request url: %s, method: %s\n", r.URL.Path, r.Method)
		body, _ := io.ReadAll(r.Body)
		NewResponse(200, "ok", map[string]string{
			"url":       r.URL.Path,
			"method":    r.Method,
			"body_size": fmt.Sprintf("%d", r.ContentLength),
			"body":      string(body),
		}).Write(w)
	})

	http.HandleFunc("/api/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request url: %s, method: %s\n", r.URL.Path, r.Method)
		w.Write([]byte("pong"))
	})

	http.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request url: %s, method: %s\n", r.URL.Path, r.Method)
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/api/v1/mini-upload", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request url: %s, method: %s\n", r.URL.Path, r.Method)

		err := r.ParseMultipartForm(200 << 20)
		log.Printf("parse form err: %v\n", err)

		file, finfo, err := r.FormFile("file")
		if err != nil {
			log.Printf("parse form failed, err: %v\n", err)
			w.Write([]byte("fail to parse form"))
			return
		}
		defer file.Close()

		randStr := fmt.Sprintf(".%d", time.Now().UnixNano())
		F, err := os.Create(finfo.Filename + randStr + ".upload")
		if err != nil {
			log.Printf("create file failed, err: %v\n", err)
			w.Write([]byte("fail to create file"))
			return
		}
		defer F.Close()
		n, err := io.Copy(F, file)
		if err != nil {
			log.Printf("copy file failed, err: %v\n", err)
			w.Write([]byte("fail to copy file"))
			return
		}
		log.Printf("copy %d bytes\n", n)
		w.Write([]byte("ok"))
	})

	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
