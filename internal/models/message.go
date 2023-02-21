package models

import "net/http"

type ClientMessage struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

type ProxyResponse struct {
	ID      string      `json:"id"`
	Status  string      `json:"status"`
	Headers http.Header `json:"headers"`
	Length  int64       `json:"length"`
}
