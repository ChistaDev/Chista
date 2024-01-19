package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func set_Leak_router() (*http.Request, *httptest.ResponseRecorder, error) {
	r := gin.New()
	r.POST("/api/v1/leak", GetLeaks)

	url := "/api/v1/leak?email=resul.bozburun@owasp.org"
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return req, httptest.NewRecorder(), err
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w, nil

}

func TestLeakController(t *testing.T) {
	a := assert.New(t)

	// Initiliaze WS connection to test
	initWsConnection()

	// Request to mock router
	req, w, err := set_Leak_router()
	if err != nil {
		a.Error(err)
	}

	a.Equal(http.MethodPost, req.Method, "HTTP Req error. Unexpected method!")
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		a.Error(err)
	}

	fmt.Printf("Response body from Leak endpoint: %v", string(body))

}
