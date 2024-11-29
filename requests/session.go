package requests

import (
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
)

type Interface interface {
	Post(url string, args ...any) *Request
	Get(url string, args ...any) *Request
	Delete(url string, args ...any) *Request
	Put(url string, args ...any) *Request
	Patch(url string, args ...any) *Request
	Options(url string, args ...any) *Request
	Head(url string, args ...any) *Request
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

func (s *Session) Post(url string, args ...any) *Request {
	return Post(url, args...).WithClient(s.client)
}

func (s *Session) Get(url string, args ...any) *Request {
	return Get(url, args...).WithClient(s.client)
}

func (s *Session) Delete(url string, args ...any) *Request {
	return Delete(url, args...).WithClient(s.client)
}

func (s *Session) Put(url string, args ...any) *Request {
	return Put(url, args...).WithClient(s.client)
}

func (s *Session) Patch(url string, args ...any) *Request {
	return Patch(url, args...).WithClient(s.client)
}

func (s *Session) Options(url string, args ...any) *Request {
	return Patch(url, args...).WithClient(s.client)
}

func (s *Session) Head(url string, args ...any) *Request {
	return Head(url, args...).WithClient(s.client)
}
