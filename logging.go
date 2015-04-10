package accesslog

import (
	//"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"strings"
)

// Logger defines needed interfaces
// One can use any logger with this interface
type Logger interface {
	Printf(format string, v ...interface{})
}

type Log struct {
	Logger Logger
}

// Wrap ResponseWriter for determination of status and content length
type LogResponseWriter struct {
	http.ResponseWriter
	status        int
	writtenLength int
}

func (lrw *LogResponseWriter) WriteHeader(status int) {
	lrw.status = status
	lrw.ResponseWriter.WriteHeader(status)
}

func (lrw *LogResponseWriter) Write(b []byte) (int, error) {
	written, err := lrw.ResponseWriter.Write(b)
	lrw.writtenLength += written
	return written, err
}

func New() *Log {
	l := &Log{
		Logger: log.New(os.Stdout, "", 0),
	}
	return l
}

func (l *Log) logging(rw *LogResponseWriter, r *http.Request) {
	endTime := time.Now()
	userAgent := r.UserAgent()
	referer := r.Referer()
	if referer == "" {
		referer = "-"
	}
	if userAgent == "" {
		userAgent = "-"
	}
	
	ip := strings.Split(r.RemoteAddr, ":")[0]
	
	// IP user-identifier user-id [datetime] "method url protocol_version" status length "referer" "user-agent"
	const format = "%s - - [%s] \"%s %s %s\" %d %d \"%s\" \"%s\""
	// "%d/%b/%Y:%H:%M:%S %z" 
	const layout = "2/Jan/2006:15:04:05 -0700"
	l.Logger.Printf(format,
		ip,
		endTime.Format(layout),
		r.Method,
		r.URL.String(),
		r.Proto,
		rw.status,
		rw.writtenLength,
		referer,
		userAgent)
}

// Middleware implementing fun(next http.Handler) http.Handler
//
// router := r2router.NewSeeforRouter()
// logger := accesslog.New()
// router.Before(logger.Handler)
func (l *Log) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		lwr := &LogResponseWriter{rw, 0, 0}
		defer l.logging(lwr, r)
		next.ServeHTTP(lwr, r)
	})
}

// Middleware implementing func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
//
// n := negroni.New()
// logger := accesslog.New()
// n.UseFunc(logger.HandlerFuncWithNext)
func (l *Log) HandlerFuncWithNext(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	lwr := &LogResponseWriter{rw, 0, 0}
	defer l.logging(lwr, r)
	next(lwr, r)
}
