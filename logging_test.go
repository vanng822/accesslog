package accesslog

import (
	"github.com/stretchr/testify/assert"
	"github.com/vanng822/r2router"
	"github.com/codegangsta/negroni"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	router := r2router.NewSeeforRouter()
	log := NewLog()
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

func TestSeeforRecoveryPrintStack(t *testing.T) {
	router := r2router.NewRouter()
	n := negroni.New()
	log := NewLog()
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
