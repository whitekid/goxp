package requests

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	*http.Response
}

// Success return error if response failed
func (r *Response) Success() error {
	if http.StatusOK <= r.StatusCode && r.StatusCode < http.StatusMultipleChoices {
		return nil
	}

	return fmt.Errorf("request failed with status %d: %s", r.StatusCode, r.Status)
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
// Depreciated: please use goxp.ReadJSON()
func (r *Response) JSON(v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// XML decode response body to xml
// caller should close body
// Depreciated: please use goxp.ReadXML()
func (r *Response) XML(v any) error {
	return xml.NewDecoder(r.Body).Decode(v)
}
