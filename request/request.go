package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/whitekid/goxp/fx"
	"github.com/whitekid/goxp/log"
)

const (
	HeaderUserAgent     = "User-Agent"
	HeaderReferer       = "Referer"
	HeaderContentType   = "Content-Type"
	HeaderAuthorization = "Authorization"
	HeaderLocation      = "Location"
)

var (
	MIMEApplicationForm = "application/x-www-form-urlencoded"

	MIMEApplicationJSON = mimeByExt(".json")
	MIMEImagePNG        = mimeByExt(".png")
	MIMEImageJPEG       = mimeByExt(".jpg")
	MIMEImageGIF        = mimeByExt(".gif")
	MIMEVCard           = mimeByExt(".vcf")
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
	followRedirect    bool
	options           []option
	client            *http.Client
}

func Post(url string, args ...interface{}) *Request   { return New(http.MethodPost, url, args...) }
func Get(url string, args ...interface{}) *Request    { return New(http.MethodGet, url, args...) }
func Delete(url string, args ...interface{}) *Request { return New(http.MethodDelete, url, args...) }
func Put(url string, args ...interface{}) *Request    { return New(http.MethodPut, url, args...) }
func Patch(url string, args ...interface{}) *Request  { return New(http.MethodPatch, url, args...) }
func Head(url string, args ...interface{}) *Request   { return New(http.MethodHead, url, args...) }

// New create new request
func New(method, URL string, args ...interface{}) *Request {
	if len(args) > 0 {
		URL = fmt.Sprintf(URL, args...)
	}

	return &Request{
		method:     method,
		URL:        URL,
		header:     http.Header{},
		query:      url.Values{},
		formValues: url.Values{},
		jsonValues: make([]interface{}, 0),
	}
}

type option interface {
	apply()
}

type funcOption struct {
	fn func()
}

func (opt *funcOption) apply() { opt.fn() }

func newFuncOption(fn func()) option {
	return &funcOption{
		fn: fn,
	}
}

func (r *Request) addOpt(opt option) *Request { r.options = append(r.options, opt); return r }
func (r *Request) addOptF(fn func()) *Request { r.addOpt(newFuncOption(fn)); return r }

func (r *Request) FollowRedirect(follow bool) *Request {
	return r.addOptF(func() { r.followRedirect = follow })
}

func (r *Request) Header(key, value string) *Request {
	return r.addOptF(func() { r.header.Add(key, value) })
}

func (r *Request) Headers(headers map[string]string) *Request {
	fx.ForEachMap(headers, func(k string, v string) { r.addOptF(func() { r.header.Add(k, v) }) })
	return r
}

func (r *Request) ContentType(contentType string) *Request {
	return r.addOptF(func() { r.header.Set(HeaderContentType, contentType) })
}

func (r *Request) AuthBasic(user, password string) *Request {
	return r.addOptF(func() {
		r.basicAuthUser = user
		r.basicAuthPassword = password
	})
}

func (r *Request) AuthBearer(token string) *Request {
	return r.addOptF(func() { r.Header(HeaderAuthorization, "Bearer "+token) })
}

func (r *Request) AuthToken(token string) *Request {
	return r.addOptF(func() { r.Header(HeaderAuthorization, "Token "+token) })
}

// Query set query parameters
func (r *Request) Query(key, value string) *Request {
	return r.addOptF(func() { r.query.Add(key, value) })
}

func (r *Request) Queries(params map[string]string) *Request {
	fx.ForEachMap(params, func(k string, v string) { r.addOptF(func() { r.query.Add(k, v) }) })
	return r
}

func (r *Request) Form(key, value string) *Request {
	return r.addOptF(func() { r.formValues.Add(key, value) })
}

func (r *Request) Forms(values map[string]string) *Request {
	fx.ForEachMap(values, func(k string, v string) { r.addOptF(func() { r.formValues.Add(k, v) }) })
	return r
}

func (r *Request) JSON(value interface{}) *Request {
	return r.addOptF(func() { r.jsonValues = append(r.jsonValues, value) })
}

// Body set request body
func (r *Request) Body(body io.Reader) *Request {
	return r.addOptF(func() { r.body = body })
}

func (r *Request) WithClient(client *http.Client) *Request {
	return r.addOptF(func() { r.client = client })
}

func (r *Request) makeRequest() (*http.Request, error) {
	for _, opt := range r.options {
		opt.apply()
	}

	u := r.URL
	if len(r.query) > 0 {
		URL, err := url.Parse(u)
		if err != nil {
			return nil, err
		}

		URL.RawQuery = url.Values(fx.MergeMap(URL.Query(), r.query)).Encode()
		u = URL.String()
	}

	var body io.Reader

	switch r.method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		switch {
		case len(r.formValues) > 0:
			r.header.Set(HeaderContentType, MIMEApplicationForm)
			body = strings.NewReader(r.formValues.Encode())
		case len(r.jsonValues) > 0:
			r.header.Set(HeaderContentType, MIMEApplicationJSON)

			buffer := fx.Map(r.jsonValues, func(v interface{}) io.Reader {
				buf := &bytes.Buffer{}
				if err := json.NewEncoder(buf).Encode(v); err != nil {
					log.Errorf("encode error: %v", err)
				}
				return buf
			})
			body = io.MultiReader(buffer...)

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

	fx.ForEachMap(r.header, func(k string, headers []string) {
		fx.ForEach(headers, func(i int, v string) { req.Header.Add(k, v) })
	})

	return req, nil
}

// Do call http request
func (r *Request) Do(ctx context.Context) (*Response, error) {
	req, err := r.makeRequest()
	if err != nil {
		return nil, err
	}

	var client *http.Client
	if r.client == nil {
		client = &http.Client{}
	} else {
		client = r.client
	}

	if !r.followRedirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }
	}

	resp, err := client.Do(req.WithContext(ctx))
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

// String return body as string
//
// string() read all data from response.Body. So that if you call more time, it returns empty string
func (r *Response) String() string {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	return string(body)
}

// JSON decode response body to json
// caller should close body
func (r *Response) JSON(v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
