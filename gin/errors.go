package gin

import (
	"fmt"
	"net/http"

	"github.com/juju/errors"
	"github.com/whitekid/go-utils/log"
	"go.uber.org/zap"
)

var (
	loggerErr = log.WithOptions(zap.AddCallerSkip(1))
)

type errorResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

func newErrorResponse(err error) *errorResponse {
	resp := errorResponse{
		Code:    getStatus(err),
		Message: err.Error(),
		Detail:  fmt.Sprintf("%+v", err),
	}

	switch {
	case resp.Code >= 500:
		loggerErr.Debugf("Code: %d, Error: %+v", resp.Code, err)
	case resp.Code >= 400:
		loggerErr.Debugf("Code: %d, Error: %s", resp.Code, err)
	}

	return &resp
}

func getStatus(err error) int {
	switch {
	case errors.IsUnauthorized(err):
		return http.StatusUnauthorized
	case errors.IsBadRequest(err):
		return http.StatusBadRequest
	case errors.IsForbidden(err):
		return http.StatusForbidden
	case errors.IsNotFound(err):
		return http.StatusNotFound
	case errors.IsMethodNotAllowed(err):
		return http.StatusMethodNotAllowed
	case errors.IsAlreadyExists(err):
		return http.StatusConflict
	case errors.IsNotImplemented(err):
		return http.StatusNotImplemented
	}

	return http.StatusInternalServerError
}
