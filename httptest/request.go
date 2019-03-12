package httptest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

// Options for request
type Options struct {
	Headers map[string]string
	Params  map[string]string
	Forms   url.Values
	Body    io.Reader   // for raw data
	JSON    interface{} // for json request
}

type ResponseRecorder struct {
	*httptest.ResponseRecorder
}

// Request send test request for http.Handler
func Request(handler http.Handler, method, path string, opts Options) (w *ResponseRecorder) {
	if method == http.MethodGet && len(opts.Params) > 0 {
		values := url.Values{}
		for k, v := range opts.Params {
			values[k] = []string{v}
		}
		path = path + "?" + values.Encode()
	}

	// post body
	var body io.Reader
	var contentType string
	if method == http.MethodPost {
		if opts.Body != nil {
			body = opts.Body
		}

		if len(opts.Forms) > 0 {
			body = strings.NewReader(opts.Forms.Encode())
			contentType = "application/x-www-form-urlencoded"
		}

		if opts.JSON != nil {
			buf := &bytes.Buffer{}
			json.NewEncoder(buf).Encode(opts.JSON)

			body = buf
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

	w = &ResponseRecorder{httptest.NewRecorder()}
	handler.ServeHTTP(w, req)

	return
}

func (r *ResponseRecorder) CloseNotify() <-chan bool {
	return nil
}

func (r *ResponseRecorder) JSON(intf interface{}) error {
	return json.NewDecoder(r.Body).Decode(intf)
}
