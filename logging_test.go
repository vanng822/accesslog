package accesslog

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/stretchr/testify/assert"
	"github.com/vanng822/r2router"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	router := r2router.NewSeeforRouter()
	log := New()
	router.Before(log.Handler)

	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, http.StatusText(http.StatusInternalServerError))
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	// get
	res, err := http.Get(ts.URL + "/user/keys/testing")
	assert.Nil(t, err)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, string(content), "Internal Server Error")
}

func TestHandlerFuncWithNext(t *testing.T) {
	router := r2router.NewRouter()
	n := negroni.New()
	log := New()
	n.UseFunc(log.HandlerFuncWithNext)

	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		fmt.Fprint(w, p.Get("id"))
	})
	n.UseHandler(router)
	ts := httptest.NewServer(n)
	defer ts.Close()

	// get
	res, err := http.Get(ts.URL + "/user/keys/testing")
	assert.Nil(t, err)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Contains(t, string(content), "testing")
}

type mylog struct {
	output string
}

func (m *mylog) Info(v string) {
	m.output = fmt.Sprint(v)
}

func TestLogOutput(t *testing.T) {
	router := r2router.NewSeeforRouter()
	log := New()
	m := &mylog{}
	log.Logger = m

	router.Before(log.Handler)

	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, http.StatusText(http.StatusInternalServerError))
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	// get
	res, err := http.Get(ts.URL + "/user/keys/testing")
	assert.Nil(t, err)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, string(content), "Internal Server Error")
	assert.Regexp(t, "127.0.0.1 - - \\[.+\\] \"GET /user/keys/testing HTTP/1.1\" 500 21 \"-\" \"Go [.\\d]+ package http", m.output)
}

func TestLoggerFuncOutput(t *testing.T) {
	router := r2router.NewSeeforRouter()
	log := New()
	var output string
	log.Logger = LoggerFunc(func(v string) {
		output = fmt.Sprint(v)
	})

	router.Before(log.Handler)

	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, http.StatusText(http.StatusInternalServerError))
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	// get
	res, err := http.Get(ts.URL + "/user/keys/testing")
	assert.Nil(t, err)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, string(content), "Internal Server Error")
	assert.Regexp(t, "127.0.0.1 - - \\[.+\\] \"GET /user/keys/testing HTTP/1.1\" 500 21 \"-\" \"Go [.\\d]+ package http", output)
}

// syslog func interface
func TestLoggerFuncSyslogOutput(t *testing.T) {
	router := r2router.NewSeeforRouter()
	log := New()
	var output string
	log.Logger = WrapSyslog(func(v string) (err error) {
		output = fmt.Sprint(v)
		return nil
	})

	router.Before(log.Handler)

	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, http.StatusText(http.StatusInternalServerError))
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	// get
	res, err := http.Get(ts.URL + "/user/keys/testing")
	assert.Nil(t, err)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, string(content), "Internal Server Error")
	assert.Regexp(t, "127.0.0.1 - - \\[.+\\] \"GET /user/keys/testing HTTP/1.1\" 500 21 \"-\" \"Go [.\\d]+ package http", output)
}
