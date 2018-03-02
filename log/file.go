package log

import (
	"github.com/apex/log"
	"os"
	"fmt"
	"strings"
)

// FileHandler implements a multi handler system for logging to files
type FileHandler struct {
	Path string
}

// NewFile creates a new multi handler system for logging
func NewFile(path string) *FileHandler {
	return &FileHandler{
		Path: path,
	}
}

// HandleLog implements the log.Handler interface
func (h *FileHandler) HandleLog(e *log.Entry) error {
	d := e.Timestamp.Format("2006-01-02")
	path := strings.Replace(h.Path, "%date%", d, -1)
	f, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	t := e.Timestamp.Format("15:04:05")
	_, err = fmt.Fprintf(f, "%s %-25s", t, e.Message)
	if err != nil {
		return err
	}

	for k, v := range e.Fields {
		_, err = fmt.Fprintf(f, " %s=%s", k, v)
		if err != nil {
			return err
		}
	}
	fmt.Fprintln(f)
	return nil
}