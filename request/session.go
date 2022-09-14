package request

import (
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
)

type Interface interface {
	Post(url string, args ...interface{}) *Request
	Get(url string, args ...interface{}) *Request
	Delete(url string, args ...interface{}) *Request
	Put(url string, args ...interface{}) *Request
	Patch(url string, args ...interface{}) *Request
	Head(url string, args ...interface{}) *Request
}

func NewSession(client *http.Client) Interface {
	if client == nil {
		jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		client = &http.Client{Jar: jar}
	}

	return &Session{
		client: client,
	}
}

// Session shares same clients
type Session struct {
	client *http.Client
}

func (s *Session) Post(url string, args ...interface{}) *Request {
	return Post(url, args...).WithClient(s.client)
}

func (s *Session) Get(url string, args ...interface{}) *Request {
	return Get(url, args...).WithClient(s.client)
}

func (s *Session) Delete(url string, args ...interface{}) *Request {
	return Delete(url, args...).WithClient(s.client)
}

func (s *Session) Put(url string, args ...interface{}) *Request {
	return Put(url, args...).WithClient(s.client)
}

func (s *Session) Patch(url string, args ...interface{}) *Request {
	return Patch(url, args...).WithClient(s.client)
}

func (s *Session) Head(url string, args ...interface{}) *Request {
	return Head(url, args...).WithClient(s.client)
}
