package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Printf("Run main")
	router := mux.NewRouter()
	log.Fatal(http.ListenAndServe(":"+"8080", router))
}
