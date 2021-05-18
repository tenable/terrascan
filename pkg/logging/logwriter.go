package logging

import (
	gorillaHandlers "github.com/gorilla/handlers"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// LogWriter a type that can write logs to a zap.SugaredLogger object
type LogWriter struct {
	logger *zap.SugaredLogger
}

// GetLogWriter creates and fetches a LogWriter object
func GetLogWriter(logger *zap.SugaredLogger) LogWriter {
	return LogWriter{
		logger: logger,
	}
}

// Write writes the byte array to the logger object
func (l LogWriter) Write(p []byte) (n int, err error) {
	l.logger.Info(string(p))
	return len(p), nil
}

// buildCommonLogLine builds a log entry for req in Apache Common Log Format.
// ts is the timestamp with which the entry should be logged.
// status and size are used to provide the response HTTP status and size.
func buildCommonLogLine(req *http.Request, url url.URL, ts time.Time, status int, size int) []byte {
	username := "-"
	if url.User != nil {
		if name := url.User.Username(); name != "" {
			username = name
		}
	}

	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		host = req.RemoteAddr
	}

	uri := req.RequestURI

	// Requests using the CONNECT method over HTTP/2.0 must use
	// the authority field (aka r.Host) to identify the target.
	// Refer: https://httpwg.github.io/specs/rfc7540.html#CONNECT
	if req.ProtoMajor == 2 && req.Method == "CONNECT" {
		uri = req.Host
	}
	if uri == "" {
		uri = url.RequestURI()
	}

	buf := make([]byte, 0, 3*(len(host)+len(username)+len(req.Method)+len(uri)+len(req.Proto)+50)/2)
	buf = append(buf, host...)
	buf = append(buf, " - "...)
	buf = append(buf, username...)
	buf = append(buf, ` "`...)
	buf = append(buf, req.Method...)
	buf = append(buf, " "...)
	buf = append(buf, uri...)
	buf = append(buf, " "...)
	buf = append(buf, req.Proto...)
	buf = append(buf, `" `...)
	buf = append(buf, strconv.Itoa(status)...)
	buf = append(buf, " "...)
	buf = append(buf, strconv.Itoa(size)...)
	return buf
}

// WriteRequestLog logs an http(s) request to a write object
func WriteRequestLog(writer io.Writer, params gorillaHandlers.LogFormatterParams) {
	writer.Write(buildCommonLogLine(params.Request, params.URL, params.TimeStamp, params.StatusCode, params.Size))
}
