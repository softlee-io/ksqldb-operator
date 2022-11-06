package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leewoobin789/test-camunda/producer-service/src/controller/endpoint"
)

var port string = "8080"

func main() {
	router := mux.NewRouter()

	handlers := endpoint.ReturnBundle()
	for _, h := range handlers {
		path := h.GetInfo().Path
		method := string(h.GetInfo().Method)
		router.HandleFunc(path, h.Run).Methods(method)
	}

	if err := http.ListenAndServe(":"+port, router); err != nil {
		panic(err)
	}
}
