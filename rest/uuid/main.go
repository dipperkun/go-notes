package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
)

type uuid struct {
}

func (u *uuid) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		genUUID(w, r)
		return
	}
	http.NotFound(w, r)
	return
}

func genUUID(w http.ResponseWriter, r *http.Request) {
	buf := make([]byte, 10)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, fmt.Sprintf("%x", buf))
}

func main() {
	mux := &uuid{}
	log.Fatal(http.ListenAndServe(":8000", mux))
}
