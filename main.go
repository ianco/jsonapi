package main

import (
	"log"
	"net/http"
	rest "github.com/ianco/jsonapi/rest"
)

func main() {

	router := rest.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
