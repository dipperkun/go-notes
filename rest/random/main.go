package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/float", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, rand.Float64())
	})
	mux.HandleFunc("/int", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, rand.Intn(1000_000))
	})

	log.Fatal(http.ListenAndServe(":8000", mux))
}
