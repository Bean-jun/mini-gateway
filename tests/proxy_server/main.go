package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
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

	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
