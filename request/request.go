package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	ContentTypeJSON = "application/json; charset=UTF-8"
	ContentTypeForm = "application/x-www-form-urlencoded; charset=UTF-8"

	headerContentType = "Content-Type"
)

type Request struct {
	URL               string
	method            string
	header            http.Header
	basicAuthUser     string
	basicAuthPassword string
	query             url.Values
	formValues        url.Values
	jsonValues        []interface{}
	body              io.Reader
	client            *http.Client
}

func Post(url string, args ...interface{}) *Request   { return New(http.MethodPost, url, args...) }
func Get(url string, args ...interface{}) *Request    { return New(http.MethodGet, url, args...) }
func Delete(url string, args ...interface{}) *Request { return New(http.MethodDelete, url, args...) }
func Put(url string, args ...interface{}) *Request    { return New(http.MethodPut, url, args...) }
func Patch(url string, args ...interface{}) *Request  { return New(http.MethodPatch, url, args...) }

// New create new request
func New(method, u string, args ...interface{}) *Request {
	if len(args) > 0 {
		u = fmt.Sprintf(u, args...)
	}

	return &Request{
		method:     method,
		URL:        u,
		header:     http.Header{},
		query:      url.Values{},
		formValues: url.Values{},
		jsonValues: make([]interface{}, 0),
	}
}

func (r *Request) FollowRedirect(follow bool) *Request {
	if r.client == nil || r.client == http.DefaultClient {
		r.client = &http.Client{}
	}

	if follow {
		r.client.CheckRedirect = nil
	} else {
		r.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return r
}

func (r *Request) Header(key, value string) *Request {
	r.header.Add(key, value)
	return r
}

func (r *Request) Headers(headers map[string]string) *Request {
	for k, v := range headers {
		r.Header(k, v)
	}

	return r
}

func (r *Request) ContentType(contentType string) *Request {
	r.header.Set(headerContentType, contentType)
	return r
}

func (r *Request) AuthBasic(user, password string) *Request {
	r.basicAuthUser = user
	r.basicAuthPassword = password
	return r
}

func (r *Request) AuthBearer(token string) *Request {
	return r.Header("Authorization", "Bearer "+token)
}

func (r *Request) AuthToken(token string) *Request {
	return r.Header("Authorization", "Token "+token)
}

// Query set query parameters
func (r *Request) Query(key, value string) *Request {
	r.query.Add(key, value)
	return r
}

func (r *Request) Queries(params map[string]string) *Request {
	for k, v := range params {
		r.Query(k, v)
	}
	return r
}

func (r *Request) Form(key, value string) *Request {
	r.formValues.Add(key, value)

	return r
}

func (r *Request) Forms(values map[string]string) *Request {
	for k, v := range values {
		r.Form(k, v)
	}
	return r
}

func (r *Request) JSON(value interface{}) *Request {
	r.jsonValues = append(r.jsonValues, value)
	return r
}

// Body set request body
func (r *Request) Body(body io.Reader) *Request {
	r.body = body
	return r
}

func (r *Request) WithClient(client *http.Client) *Request { r.client = client; return r }

func (r *Request) makeRequest() (*http.Request, error) {
	u := r.URL
	if len(r.query) > 0 {
		URL, err := url.Parse(u)
		if err != nil {
			return nil, err
		}

		params := url.Values{}
		query := URL.Query()
		if len(query) > 0 {
			for k, v := range URL.Query() {
				params[k] = v
			}
		}

		for k, v := range r.query {
			params[k] = v
		}
		URL.RawQuery = params.Encode()

		u = URL.String()
	}

	var body io.Reader

	switch r.method {
	case http.MethodPost, http.MethodPut:
		switch {
		case len(r.formValues) > 0:
			r.header.Set(headerContentType, ContentTypeForm)
			body = strings.NewReader(r.formValues.Encode())
		case len(r.jsonValues) > 0:
			r.header.Set(headerContentType, ContentTypeJSON)
			buffer := &bytes.Buffer{}
			for _, js := range r.jsonValues {
				buf := &bytes.Buffer{}
				if err := json.NewEncoder(buf).Encode(js); err == nil {
					buf.WriteTo(buffer)
				}
			}
			body = buffer
		case r.body != nil:
			body = r.body
		}
	}

	req, err := http.NewRequest(r.method, u, body)
	if err != nil {
		return nil, err
	}

	if r.basicAuthUser != "" {
		req.SetBasicAuth(r.basicAuthUser, r.basicAuthPassword)
	}

	for k, headers := range r.header {
		for _, h := range headers {
			req.Header.Add(k, h)
		}
	}

	return req, nil
}

// Do call http request
func (r *Request) Do(ctx context.Context) (*Response, error) {
	req, err := r.makeRequest()
	if err != nil {
		return nil, err
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	client := http.DefaultClient
	if r.client != nil {
		client = r.client
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return &Response{resp}, nil
}

type Response struct {
	*http.Response
}

func (r *Response) Success() bool {
	return http.StatusOK <= r.StatusCode && r.StatusCode < http.StatusMultipleChoices
}

func (r *Response) String() string {
	data, _ := ioutil.ReadAll(r.Body)
	return string(data)
}

// JSON decode response body to json
// caller should close body
func (r *Response) JSON(v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
