## Accesslog

Accesslog is a middleware with implementing func(next http.Handler) http.Handler. See http://en.wikipedia.org/wiki/Common_Log_Format

## Example

```go	
package main

import (
	"fmt"
	"github.com/vanng822/accesslog"
	"github.com/vanng822/r2router"
	"net/http"
	"log"
)

func main() {
	seefor := r2router.NewSeeforRouter()
	l := accesslog.NewLog()
	seefor.Before(l.Handler)
	seefor.Get("/hello/:name", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		fmt.Fprintf(w, "Hello %s!", p.Get("name"))
	})
	log.Fatal(http.ListenAndServe(":8080", seefor))
}
```	
