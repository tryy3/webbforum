package log

import (
	"github.com/apex/log"
	"github.com/hashicorp/go-multierror"
)

// MultiHandler implements a multi handler system for logging to multiple handlers
type MultiHandler struct {
	Handlers []log.Handler
}

// NewMulti creates a new multi handler system for logging to multiple handlers
func NewMulti(h ...log.Handler) *MultiHandler {
	return &MultiHandler{
		Handlers: h,
	}
}

// HandleLog will redirect a log entry to multiple loggers
func (h *MultiHandler) HandleLog(e *log.Entry) error {
	var result *multierror.Error
	for _, handler := range h.Handlers {
		if err := handler.HandleLog(e); err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result.ErrorOrNil()
}
