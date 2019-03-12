package request

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	ContentTypeJSON = "application/json"
	ContentTypeForm = "application/x-www-form-urlencoded"

	headerContentType = "Content-Type"
)

var DefaultTimeout time.Duration = 5 * time.Second

type Request struct {
	URL        string
	method     string
	header     http.Header
	params     url.Values
	formValues url.Values
	jsonValues []interface{}
}

var (
	client *http.Client
)

func init() {
	client = &http.Client{
		Timeout: DefaultTimeout,
	}
}

func Post(url string) *Request {
	return New(http.MethodPost, url)
}

func Get(url string) *Request {
	return New(http.MethodGet, url)
}

func Delete(url string) *Request {
	return New(http.MethodDelete, url)
}

func Put(url string) *Request {
	return New(http.MethodPut, url)
}

func Patch(url string) *Request {
	return New(http.MethodPatch, url)
}

//New create new request
func New(method, u string) *Request {
	return &Request{
		method:     method,
		URL:        u,
		header:     http.Header{},
		params:     url.Values{},
		formValues: url.Values{},
		jsonValues: make([]interface{}, 0),
	}
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

func (r *Request) AuthBearer(token string) *Request {
	return r.Header("Authorization", "Bearer "+token)
}

func (r *Request) AuthToken(token string) *Request {
	return r.Header("Authorization", "Token "+token)
}

func (r *Request) Param(key, value string) *Request {
	r.params.Add(key, value)
	return r
}

func (r *Request) Params(params map[string]string) *Request {
	for k, v := range params {
		r.Param(k, v)
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

func (r *Request) makeRequest() (*http.Request, error) {
	u := r.URL
	if len(r.params) > 0 {
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

		for k, v := range r.params {
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
		}
	}

	req, err := http.NewRequest(r.method, u, body)
	if err != nil {
		return nil, err
	}

	for k, headers := range r.header {
		for _, h := range headers {
			req.Header.Add(k, h)
		}
	}

	return req, nil
}

func (r *Request) Do() (*Response, error) {
	req, err := r.makeRequest()
	if err != nil {
		return nil, err
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
	return 200 <= r.StatusCode && r.StatusCode < 300
}

func (r *Response) String() string {
	data, _ := ioutil.ReadAll(r.Body)
	return string(data)
}

func (r *Response) JSON(v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}