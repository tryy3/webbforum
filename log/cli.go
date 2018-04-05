package log

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

// start time.
var start = time.Now()

var bold = color.New(color.Bold)

// Colors mapping.
var Colors = [...]*color.Color{
	log.DebugLevel: color.New(color.FgHiCyan),
	log.InfoLevel:  color.New(color.FgGreen),
	log.WarnLevel:  color.New(color.FgYellow),
	log.ErrorLevel: color.New(color.FgRed),
	log.FatalLevel: color.New(color.FgRed),
}

// Strings mapping.
var Strings = [...]string{
	log.DebugLevel: "^",
	log.InfoLevel:  "*",
	log.WarnLevel:  "*",
	log.ErrorLevel: "-",
	log.FatalLevel: "-",
}

// CliHandler writes all the logging to the terminal/console
type CliHandler struct {
	mu      sync.Mutex
	Writer  io.Writer
	Padding int
}

// NewCli handler.
func NewCli(w io.Writer) *CliHandler {
	if f, ok := w.(*os.File); ok {
		return &CliHandler{
			Writer:  colorable.NewColorable(f),
			Padding: 3,
		}
	}

	return &CliHandler{
		Writer:  w,
		Padding: 3,
	}
}

// HandleLog takes care of formatting, parsing and outputting the log message to the terminal/console
func (h *CliHandler) HandleLog(e *log.Entry) error {
	color := Colors[e.Level]
	level := Strings[e.Level]
	names := e.Fields.Names()

	h.mu.Lock()
	defer h.mu.Unlock()

	color.Fprintf(h.Writer, "%s %-25s",
		bold.Sprintf("%s %*s", e.Timestamp.Format("15:04:05"), h.Padding+1, level),
		e.Message)

	for _, name := range names {
		if name == "source" {
			continue
		}
		fmt.Fprintf(h.Writer, " %s=%s", color.Sprint(name), e.Fields.Get(name))
	}

	fmt.Fprintln(h.Writer)

	return nil
}
