package ctx

import (
	"context"
	"fmt"
	"net/http"
)

func Serve(store Storer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := store.Fetch(r.Context())
		if err != nil {
			return
		}
		fmt.Fprint(w, data)
	}
}

type Storer interface {
	Fetch(context.Context) (string, error)
}
