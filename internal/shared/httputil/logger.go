package httputil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Logger struct {
	log *log.Logger
}

func NewLogger(log *log.Logger) *Logger {
	return &Logger{log: log}
}

type record struct {
	Time     time.Time     `json:"time"`
	Level    string        `json:"level"`
	Msg      string        `json:"msg"`
	Method   string        `json:"method"`
	Path     string        `json:"path"`
	Status   int           `json:"status"`
	Duration time.Duration `json:"duration"`
	Request  string        `json:"request,omitempty"`
	Err      string        `json:"error,omitempty"`
	IP       string        `json:"ip"`
}

func (l *Logger) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()
		lw := &logWriter{ResponseWriter: w}

		requestBody, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewReader(requestBody))

		next(lw, r)

		rec := record{
			Time:     start,
			Msg:      "http_request",
			Method:   r.Method,
			Path:     r.URL.Path,
			IP:       r.RemoteAddr,
			Status:   lw.status,
			Duration: time.Since(start),
		}

		if rec.Status < 400 {
			rec.Level = "INFO"
		} else {
			rec.Level = "ERROR"
			if rec.Status < 500 {
				rec.Level = "WARN"
			}
			rec.Request = string(requestBody)
			rec.Err = lw.err
		}

		logMarshaled, _ := json.Marshal(rec)
		l.log.Print(string(logMarshaled))
	}
}

type logWriter struct {
	http.ResponseWriter
	status int
	err    string
}

func (lw *logWriter) WriteHeader(statusCode int) {
	lw.status = statusCode
	lw.ResponseWriter.WriteHeader(statusCode)
}

func (lw *logWriter) Write(body []byte) (int, error) {
	if lw.status == 0 {
		lw.status = http.StatusOK
	}

	if lw.status >= 400 {
		var handlerErr map[string]string
		if err := json.Unmarshal(body, &handlerErr); err != nil {
			panic("invalid handler error format (use httputil WriteErr functions)")
		}

		lw.err = handlerErr["log_err"]
		clientErr := handlerErr["client_err"]

		return lw.ResponseWriter.Write([]byte(fmt.Sprintf(`{"error":%q}`, clientErr)))
	}

	return lw.ResponseWriter.Write(body)
}
