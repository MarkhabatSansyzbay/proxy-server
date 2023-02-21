package service

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"proxyserver/internal/models"

	"github.com/google/uuid"
)

var ErrMethod = errors.New("method is not allowed")

type ReverseProxy interface {
	ResponseToMessage(request *models.ClientMessage) (models.ProxyResponse, error)
}

type Proxy struct {
	store *Store
}

func NewProxy() *Proxy {
	return &Proxy{
		store: NewStore(),
	}
}

func (p *Proxy) ResponseToMessage(msg *models.ClientMessage) (models.ProxyResponse, error) {
	if strings.ToUpper(msg.Method) != "GET" {
		return models.ProxyResponse{}, ErrMethod
	}

	clientReq := modifyClientRequest(msg)
	proxyResp, isStored := p.store.get(clientReq)
	if isStored {
		return proxyResp, nil
	}

	req, err := http.NewRequest(msg.Method, msg.URL, nil)
	if err != nil {
		return models.ProxyResponse{}, fmt.Errorf("error creating a new request: %s", err)
	}

	for key, val := range msg.Headers {
		req.Header.Set(key, val)
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return models.ProxyResponse{}, fmt.Errorf("error sending http request: %s", err)
	}

	proxyResp = p.proxyResponse(resp)
	p.store.save(clientReq, proxyResp)

	return proxyResp, nil
}

func (p *Proxy) proxyResponse(resp *http.Response) models.ProxyResponse {
	proxyResp := &models.ProxyResponse{
		ID:      uuid.New().String(),
		Status:  resp.Status,
		Headers: resp.Header,
		Length:  resp.ContentLength,
	}
	return *proxyResp
}

func modifyClientRequest(msg *models.ClientMessage) string {
	b := new(bytes.Buffer)
	for key, val := range msg.Headers {
		fmt.Fprintf(b, "%s:%s", key, val)
	}

	return fmt.Sprintf(
		"URL: %s Headers%s",
		msg.URL,
		b.String(),
	)
}
