package requests

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/whitekid/goxp/errors"
)

type Response struct {
	*http.Response
}

// Success return error if response failed
func (r *Response) Success() error {
	if http.StatusOK <= r.StatusCode && r.StatusCode < http.StatusMultipleChoices {
		return nil
	}

	return errors.Errorf(nil, "request failed with status %d: %s", r.StatusCode, r.Status)
}

// String return body as string
//
// string() read all data from response.Body. So that if you call more time, it returns empty string
func (r *Response) String() string {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	return string(body)
}

// JSON return body as json
func (r *Response) JSON(v any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
