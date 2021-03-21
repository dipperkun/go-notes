package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	logMid := handlers.LoggingHandler(os.Stdout, r)
	http.ListenAndServe(":8000", logMid)
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("index")
	w.Write([]byte("OK"))
}
