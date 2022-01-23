package concurrent

import (
	"testing"
	"net/http/httptest"
)

func TestRace(t *testing.T) {

	slowServer := makeDelayedServer(20 * time.Millisecond)

	fastServer := makeDelayedServer(0 * time.Millisecond)

	defer slowServer.Close()
	defer fastServer.Close()

	slowURL := slowServer.URL
	fastURL := fastServer.URL

	expected := fastURL
	got, _ := Race(slowURL, fastURL, 20*time.Millisecond)

	if got != expected {
		t.Errorf("got %q, expected %q", got, expected)
	}
}

func makeDelayedServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
	}))
}