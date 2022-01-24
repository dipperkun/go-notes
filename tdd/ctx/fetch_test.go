package ctx

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type StubStore struct {
	resp string
}

func (s *StubStore) Fetch() string {
	return s.resp
}

func (s *StubStore) Cancel() {

}
func TestServe(t *testing.T) {
	data := "hello, world"
	s := Serve(&SpyStore{data, t})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	s.ServeHTTP(resp, req)
	if resp.Body.String() != data {
		t.Errorf(`got "%s", expected "%s"`, resp.Body.String(), data)
	}
}

type SpyStore struct {
	resp string
	t    *testing.T
}

func (s *SpyStore) Fetch(ctx context.Context) (string, error) {
	data := make(chan string, 1)
	go func() {
		var res string
		for _, c := range s.resp {
			select {
			case <-ctx.Done():
				log.Println("spy store got cancelled")
				return
			default:
				time.Sleep(500 * time.Millisecond)
				res += string(c)
			}
		}
		data <- res
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case res := <-data:
		return res, nil
	}
}

func (s *SpyStore) assertCancelled() {
	s.t.Helper()
	// if !s.cancelled {
	// 	s.t.Error("store was not told to cancel")
	// }
}

func (s *SpyStore) assertNotCancelled() {
	s.t.Helper()
	// if s.cancelled {
	// 	s.t.Error("store was told to cancel")
	// }
}

func TestServe2(t *testing.T) {
	data := "hello, world"
	store := &SpyStore{
		resp: data,
		t:    t,
	}
	s := Serve(store)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx, cancel := context.WithCancel(req.Context())
	time.AfterFunc(5*time.Millisecond, cancel)
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)
	// if !store.cancelled {
	// 	t.Error("store was not told to cancel")
	// }
}

func TestServe3(t *testing.T) {
	t.Run("returns data from store", func(t *testing.T) {
		data := "hello, world"
		store := &SpyStore{resp: data, t: t}
		s := Serve(store)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		resp := httptest.NewRecorder()

		s.ServeHTTP(resp, req)

		if resp.Body.String() != data {
			t.Errorf(`got "%s", expected "%s"`, resp.Body.String(), data)
		}
	})
}

type SpyResponseWriter struct {
	written bool
}

func (s *SpyResponseWriter) Header() http.Header {
	s.written = true
	return nil
}

func (s *SpyResponseWriter) Write([]byte) (int, error) {
	s.written = true
	return 0, errors.New("not implemented")
}

func (s *SpyResponseWriter) WriteHeader(code int) {
	s.written = true
}

func TestServe5(t *testing.T) {
	// tells store to cancel work if request is cancelled
	data := "hello, world"
	store := &SpyStore{resp: data, t: t}
	s := Serve(store)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx, cancel := context.WithCancel(req.Context())
	time.AfterFunc(5*time.Millisecond, cancel)
	req = req.WithContext(ctx)

	resp := &SpyResponseWriter{}

	s.ServeHTTP(resp, req)

	if resp.written {
		t.Error("a response should not have been written")
	}
}
