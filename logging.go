// Package accesslog logs request in standard format http://en.wikipedia.org/wiki/Common_Log_Format using interface Logger.Printf(format string, v ...interface{})
//
//	package main
//	
//	import (
//		"fmt"
//		"github.com/vanng822/accesslog"
//		"github.com/vanng822/r2router"
//		"net/http"
//	)
//	
//	func main() {
//		seefor := r2router.NewSeeforRouter()
//		log := accesslog.NewLog()
//		seefor.Before(log.Handler)
//		seefor.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
//			fmt.Fprint(w, p.Get("id"))
//		})
//		http.ListenAndServe(":8080", seefor)
//	}
package accesslog

import (
	//"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

type Log struct {
	Logger Logger
	format string
}

type loggingResponse struct {
	http.ResponseWriter
	status        int
	writtenLength int
}

func (lrw *loggingResponse) WriteHeader(status int) {
	lrw.status = status
	lrw.ResponseWriter.WriteHeader(status)
}

func (lrw *loggingResponse) Write(b []byte) (int, error) {
	written, err := lrw.ResponseWriter.Write(b)
	lrw.writtenLength += written
	return written, err
}

func NewLog() *Log {
	l := &Log{
		Logger: log.New(os.Stdout, "", 0),
	}
	return l
}

func (l *Log) logging(rw *loggingResponse, r *http.Request) {
	endTime := time.Now()
	userAgent := r.UserAgent()
	referer := r.Referer()
	if referer == "" {
		referer = "-"
	}
	if userAgent == "" {
		userAgent = "-"
	}
	// IP user-identifier user-id [datetime] "method url protocol_version" status length "referer" "user-agent"
	const format = "%s - - [%s] \"%s %s %s\" %d %d \"%s\" \"%s\""
	// "%d/%b/%Y:%H:%M:%S %z" 
	const layout = "2/Jan/2006 15:04:05 -0700"
	l.Logger.Printf(format,
		r.RemoteAddr,
		endTime.Format(layout),
		r.Method,
		r.URL.String(),
		r.Proto,
		rw.status,
		rw.writtenLength,
		referer,
		userAgent)
}

func (l *Log) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		lwr := &loggingResponse{rw, 0, 0}
		defer l.logging(lwr, r)
		next.ServeHTTP(lwr, r)
	})
}

func (l *Log) HandlerFuncWithNext(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	lwr := &loggingResponse{rw, 0, 0}
	defer l.logging(lwr, r)
	next(lwr, r)
}
