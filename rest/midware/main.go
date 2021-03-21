package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/justinas/alice"
)

// middleware
func one(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("one: before req")
		h.ServeHTTP(w, r)
		fmt.Println("one: after resp")
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("index.....")
	w.Write([]byte("OK"))
}

func main() {
	http.Handle("/", one(http.HandlerFunc(index)))
	http.Handle("/city", check(cookie(http.HandlerFunc(f))))

	http.Handle("/alice", alice.New(check, cookie).Then(http.HandlerFunc(f)))

	log.Fatal(http.ListenAndServe(":8000", nil))
}

type city struct {
	Name string `json:"name"`
	Area uint64 `json:"area"`
}

// curl -i -H "Content-Type: application/json" -X POST \
// http://localhost:8000/city -d '{"name":"Beijing", "area":500}'

func f(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var tempCity city
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&tempCity)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()
		fmt.Printf("Got %s city with area of %d sq miles!\n", tempCity.Name, tempCity.Area)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("201 - Created"))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method Not Allowed"))
	}
}

func check(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("check: before req")
		if r.Header.Get("Content-type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("415 - Unsupported Media Type. Please send JSON"))
			return
		}
		h.ServeHTTP(w, r)
	})
}
func cookie(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		// Setting cookie to every API response
		cookie := http.Cookie{Name: "Server-Time(UTC)",
			Value: strconv.FormatInt(time.Now().Unix(), 10)}
		http.SetCookie(w, &cookie)
		log.Println("cookie: after resp")
	})
}
