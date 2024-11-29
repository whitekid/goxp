package httptest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/whitekid/goxp"
)

// Options for request
type Options struct {
	Headers     map[string]string
	BearerToken string
	Params      url.Values
	Forms       url.Values
	Body        io.Reader   // for raw data
	JSON        any // for json request, string or struct
}

// ResponseRecorder is httptest.ResponseRecorder wrapper
type ResponseRecorder struct {
	*httptest.ResponseRecorder
}

// NewOptions create new Options instance
func NewOptions() *Options {
	return &Options{
		Headers: map[string]string{},
		Params:  url.Values{},
		Forms:   url.Values{},
	}
}

// Request send test request for http.Handler
func Request(handler http.Handler, method, path string, opts Options) (w *ResponseRecorder) {
	defer goxp.Timer("%s %s", method, path)()

	if method == http.MethodGet && len(opts.Params) > 0 {
		path = path + "?" + opts.Params.Encode()
	}

	// post body
	var body io.Reader
	var contentType string
	if method == http.MethodPost || method == http.MethodPut {
		if opts.Body != nil {
			body = opts.Body
		}

		if len(opts.Forms) > 0 {
			body = strings.NewReader(opts.Forms.Encode())
			contentType = "application/x-www-form-urlencoded"
		}

		if opts.JSON != nil {
			s, ok := opts.JSON.(string)
			if ok {
				body = strings.NewReader(s)
			} else {
				buf := &bytes.Buffer{}
				json.NewEncoder(buf).Encode(opts.JSON)

				body = buf
			}
			contentType = "application/json"
		}
	}

	req, _ := http.NewRequest(method, path, body)

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	for k, v := range opts.Headers {
		req.Header.Set(k, v)
	}

	if opts.BearerToken != "" {
		req.Header.Set("authorization", "bearer "+opts.BearerToken)
	}

	w = &ResponseRecorder{httptest.NewRecorder()}
	handler.ServeHTTP(w, req)

	return
}

// OK return true if response has ok code (200~299)
func (r *ResponseRecorder) OK() bool {
	return 200 <= r.Code && r.Code < 300
}

// CloseNotify ...
func (r *ResponseRecorder) CloseNotify() <-chan bool {
	return nil
}

// JSON read json form response
func (r *ResponseRecorder) JSON(intf any) error {
	return json.NewDecoder(r.Body).Decode(intf)
}
