package requests

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"

	"github.com/whitekid/goxp"
	"github.com/whitekid/goxp/log"
	"github.com/whitekid/goxp/mapx"
	"github.com/whitekid/goxp/slicex"
)

const (
	HeaderUserAgent       = "User-Agent"
	HeaderReferer         = "Referer"
	HeaderContentType     = "Content-Type"
	HeaderContentEncoding = "Content-Encoding"
	HeaderAuthorization   = "Authorization"
	HeaderLocation        = "Location"

	headerAccept         = "Accept"
	headerAcceptEncoding = "Accept-Encoding"

	defaultAcceptEncoding = "gzip, deflate, br, zstd"
)

var (
	MIMEApplicationForm = "application/x-www-form-urlencoded"

	MIMEApplicationJSON = mimeByExt(".json")
	MIMEImagePNG        = mimeByExt(".png")
	MIMEImageJPEG       = mimeByExt(".jpg")
	MIMEImageGIF        = mimeByExt(".gif")
	MIMEVCard           = mimeByExt(".vcf")
)

var logger = log.Named("request")

type Request struct {
	URL               string
	method            string
	header            http.Header
	basicAuthUser     string
	basicAuthPassword string
	query             url.Values
	formValues        url.Values
	jsonValues        []any
	body              io.Reader
	noFollowRedirect  bool
	options           []option
	client            *http.Client
}

func Post(url string, args ...any) *Request    { return New(http.MethodPost, url, args...) }
func Get(url string, args ...any) *Request     { return New(http.MethodGet, url, args...) }
func Delete(url string, args ...any) *Request  { return New(http.MethodDelete, url, args...) }
func Put(url string, args ...any) *Request     { return New(http.MethodPut, url, args...) }
func Patch(url string, args ...any) *Request   { return New(http.MethodPatch, url, args...) }
func Options(url string, args ...any) *Request { return New(http.MethodOptions, url, args...) }
func Head(url string, args ...any) *Request    { return New(http.MethodHead, url, args...) }

// New create new request
func New(method, URL string, args ...any) *Request {
	if len(args) > 0 {
		URL = fmt.Sprintf(URL, args...)
	}

	return &Request{
		method: method,
		URL:    URL,
		header: http.Header{
			headerAcceptEncoding: []string{defaultAcceptEncoding},
		},
		query:      url.Values{},
		formValues: url.Values{},
		jsonValues: make([]any, 0),
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

// FollowRedirect default action will be follow redirect
func (r *Request) FollowRedirect(follow bool) *Request {
	return r.addOptF(func() { r.noFollowRedirect = !follow })
}

func (r *Request) Header(key, value string) *Request {
	return r.addOptF(func() { r.header.Add(key, value) })
}

func (r *Request) Headers(headers map[string]string) *Request {
	for k, v := range headers {
		r.addOptF(func() { r.header.Add(k, v) })
	}

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
	return r.Header(HeaderAuthorization, "Bearer "+token)
}

func (r *Request) AuthToken(token string) *Request {
	return r.Header(HeaderAuthorization, "Token "+token)
}

// Query set query parameters
func (r *Request) Query(key, value string) *Request {
	return r.addOptF(func() { r.query.Add(key, value) })
}

func (r *Request) Queries(params map[string]string) *Request {
	for k, v := range params {
		r.addOptF(func() { r.query.Add(k, v) })
	}

	return r
}

func (r *Request) Form(key, value string) *Request {
	return r.addOptF(func() { r.formValues.Add(key, value) })
}

func (r *Request) Forms(values map[string]string) *Request {
	for k, v := range values {
		r.addOptF(func() { r.formValues.Add(k, v) })
	}
	return r
}

func (r *Request) JSON(value any) *Request {
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

		URL.RawQuery = url.Values(mapx.Merge(URL.Query(), r.query)).Encode()
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

			buffer := slicex.Map(r.jsonValues, func(v any) io.Reader {
				buf := &bytes.Buffer{}
				if err := json.NewEncoder(buf).Encode(v); err != nil {
					logger.Errorf("encode error: %v", err)
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

	for k, vv := range r.header {
		for _, v := range vv {
			req.Header.Set(k, v)
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

	var client *http.Client
	if r.client == nil {
		client = &http.Client{}
	} else {
		client = r.client
	}

	if r.noFollowRedirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }
	}

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	var body io.ReadCloser
	switch enc := resp.Header.Get(HeaderContentEncoding); enc {
	case "gzip":
		r, err := gzip.NewReader(resp.Body)
		goxp.Must(err)
		body = newReadCloser(r, resp.Body)
	case "br":
		body = newReadCloser(brotli.NewReader(resp.Body), resp.Body)
	case "zstd":
		decoder, err := zstd.NewReader(resp.Body)
		goxp.Must(err)
		body = newReadCloser(decoder, resp.Body)
	case "deflate":
		body = newReadCloser(flate.NewReader(resp.Body), resp.Body)
	}

	if body != nil {
		resp.Body = body
	}

	return &Response{resp}, nil
}

type readCloser struct {
	io.Reader
	c io.Closer
}

func newReadCloser(r io.Reader, c io.Closer) io.ReadCloser {
	return &readCloser{Reader: r, c: c}
}

func (rc *readCloser) Close() error {
	if rc.c != nil {
		return rc.c.Close()
	}

	return nil
}
