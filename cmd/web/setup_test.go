package main

import (
	"net/http"
)

type myHandler struct {}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
