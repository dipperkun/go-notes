package ctx

import (
	"context"
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
	s := Serve(&StubStore{data})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	s.ServeHTTP(resp, req)
	if resp.Body.String() != data {
		t.Errorf(`got "%s", expected "%s"`, resp.Body.String(), data)
	}
}

type SpyStore struct {
	resp      string
	cancelled bool
	t         *testing.T
}

func (s *SpyStore) Fetch() string {
	time.Sleep(100 * time.Millisecond)
	return s.resp
}

func (s *SpyStore) Cancel() {
	s.cancelled = true
}

func (s *SpyStore) assertCancelled() {
	s.t.Helper()
	if !s.cancelled {
		s.t.Error("store was not told to cancel")
	}
}

func (s *SpyStore) assertNotCancelled() {
	s.t.Helper()
	if s.cancelled {
		s.t.Error("store was told to cancel")
	}
}

func TestServe2(t *testing.T) {
	data := "hello, world"
	store := &SpyStore{
		resp: data,
	}
	s := Serve(store)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx, cancel := context.WithCancel(req.Context())
	time.AfterFunc(5*time.Millisecond, cancel)
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)
	if !store.cancelled {
		t.Error("store was not told to cancel")
	}
}

func TestServe3(t *testing.T) {
	t.Run("returns data from store", func(t *testing.T) {
		data := "hello, world"
		store := &SpyStore{resp: data}
		s := Serve(store)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		resp := httptest.NewRecorder()

		s.ServeHTTP(resp, req)

		if resp.Body.String() != data {
			t.Errorf(`got "%s", expected "%s"`, resp.Body.String(), data)
		}

		if store.cancelled {
			t.Error("it should not have cancelled the store")
		}
	})
}
