package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func enableBasicAPI(
	router *httprouter.Router,
) {
	hello := &basicHandler{}

	router.OPTIONS("/hello", Empty)
	router.GET("/hello", hello.World)
	router.PUT("/hello", Empty)
	router.DELETE("/hello", Empty)
}

var _ httprouter.Handle = (&basicHandler{}).World

type basicHandler struct{}

func (b *basicHandler) World(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintln(w, "Hello, World")
}
