package main

import (
	"github.com/julienschmidt/httprouter"
)

func enableBasicAPI(
	router *httprouter.Router,
) {
	router.OPTIONS("/user", Empty)
	router.GET("/user", Empty)
	router.PUT("/user", Empty)
	router.DELETE("/user", Empty)
}
