package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/julienschmidt/httprouter"
)

func main() {
	r := httprouter.New()
	r.GET("/api/v1/go-version", goVersion)
	r.GET("/api/v1/files/:name", cat)
	log.Fatal(http.ListenAndServe(":8000", r))
}

func cmdOutput(cmd string, args ...string) string {
	out, _ := exec.Command(cmd, args...).Output()
	return string(out)
}

func goVersion(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	resp := cmdOutput("go", "version")
	io.WriteString(w, resp)
}

func cat(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println(params.ByName("name"))
	content := cmdOutput("/bin/cat", params.ByName("name"))
	fmt.Fprintf(w, content)
}
