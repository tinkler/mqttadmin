// implements the logger interface for chi router
package logger

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type logEntry struct {
	*LogFormatter
	request *http.Request
	buf     *bytes.Buffer
}

func (l *logEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	switch {
	case status < 200:
		l.buf.WriteString(wrapText(strconv.Itoa(status), CT_BLUE))
	case status < 300:
		l.buf.WriteString(wrapText(strconv.Itoa(status), CT_GREEN))
	case status < 400:
		l.buf.WriteString(wrapText(strconv.Itoa(status), CT_CYAN))
	case status < 500:
		l.buf.WriteString(wrapText(strconv.Itoa(status), CT_YELLOW))
	default:
		l.buf.WriteString(wrapText(strconv.Itoa(status), CT_RED))
	}

	l.buf.WriteString(wrapText(fmt.Sprintf(" %dB", bytes), CT_BLUE))

	l.buf.WriteString(" in ")
	if elapsed < 500*time.Millisecond {
		l.buf.WriteString(wrapText(
			elapsed.String(), CT_GREEN,
		))
	} else if elapsed < 5*time.Second {
		l.buf.WriteString(wrapText(
			elapsed.String(), CT_YELLOW,
		))
	} else {
		l.buf.WriteString(wrapText(
			elapsed.String(), CT_RED,
		))
	}

	if ConsoleLevel > LL_LOG {
		os.Stdout.Write(l.buf.Bytes())
	}
}

type LogFormatter struct {
	routePath map[string]string
}

func NewLogFormatter(noColor bool) *LogFormatter {
	return &LogFormatter{routePath: make(map[string]string)}
}

func (l *LogFormatter) AddRouteInfo(routePath map[string]string) {
	if ConsoleLevel == LL_DEBUG {
		for k, v := range routePath {
			l.routePath[k] = v
		}
	}
}

func (l *LogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &logEntry{
		LogFormatter: l,
		request:      r,
		buf:          &bytes.Buffer{},
	}

	if fileLine, ok := l.routePath[r.RequestURI]; ok {
		entry.buf.WriteString(
			wrapText(fmt.Sprintf("%s\n", fileLine), CT_BLACK),
		)
	}

	reqID := middleware.GetReqID(r.Context())
	if reqID != "" {
		entry.buf.WriteString(wrapText("["+reqID+"] ", CT_YELLOW))
	}
	entry.buf.WriteString(wrapText("\"", CT_CYAN))
	entry.buf.WriteString(wrapText(r.Method+" ", CT_PURPLE))

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	entry.buf.WriteString(wrapText(fmt.Sprintf("%s://%s%s %s\" ", scheme, r.Host, r.RequestURI, r.Proto), CT_CYAN))

	entry.buf.WriteString("from ")
	entry.buf.WriteString(r.RemoteAddr)
	entry.buf.WriteString(" - ")

	return entry
}

func (l *logEntry) Panic(v interface{}, stack []byte) {
	middleware.PrintPrettyStack(v)
}

// NewLogEntry creates a new LogEntry for the request.
// reference from go-chi/chi/v5/middleware/logger.go
func RequestLogger(f middleware.LogFormatter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := f.NewLogEntry(r)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				entry.Write(ww.Status(), ww.BytesWritten(), ww.Header(), time.Since(t1), nil)
			}()

			next.ServeHTTP(ww, middleware.WithLogEntry(r, entry))
		}
		return http.HandlerFunc(fn)
	}
}

func ChiLogger(f func(formatter *LogFormatter)) func(next http.Handler) http.Handler {
	color := true
	if runtime.GOOS == "windows" {
		color = false
	}
	formatter := NewLogFormatter(!color)
	if f != nil {
		f(formatter)
	}
	return RequestLogger(formatter)
}
