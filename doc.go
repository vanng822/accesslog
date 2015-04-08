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
//		log := accesslog.New()
//		seefor.Before(log.Handler)
//		seefor.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
//			fmt.Fprint(w, p.Get("id"))
//		})
//		http.ListenAndServe(":8080", seefor)
//	}
package accesslog


