package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amartin3659/VacationHomeRental/internal/config"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler
	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		// all fine nothing to do
	default:
		t.Error(fmt.Sprintf("Type mismatch: Expected http.Handler, go %T", v))
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler
	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		// all fine nothing to do
	default:
		t.Error(fmt.Sprintf("Type mismatch: Expected http.Handler, go %T", v))
	}
}

func TestValidAuth(t *testing.T) {
  var app config.AppConfig
  mux := routes(&app)
  
	req, _ := http.NewRequest("GET", "/admin/dashboard", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// rr means "request recorder" and is an initialized response recorder for http requests builtin the test
	// basically "faking" a client and to provide a valid request/response-cycle during tests
	rr := httptest.NewRecorder()

  mux.ServeHTTP(rr, req)

	// the test itself as a condition
	if rr.Code != http.StatusSeeOther {
		t.Errorf("TestInvalidAuth failed: unexpected response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

func TestInvalidAuth(t *testing.T) {
  var app config.AppConfig
  mux := routes(&app)
  
	req, _ := http.NewRequest("GET", "/admin/dashboard", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// rr means "request recorder" and is an initialized response recorder for http requests builtin the test
	// basically "faking" a client and to provide a valid request/response-cycle during tests
	rr := httptest.NewRecorder()
	session.Put(ctx, "user_id", 1)

  mux.ServeHTTP(rr, req)

	// the test itself as a condition
	if rr.Code != http.StatusOK {
		t.Errorf("TestInvalidAuth failed: unexpected response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
