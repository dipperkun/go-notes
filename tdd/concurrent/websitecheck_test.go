package concurrent

import (
	"reflect"
	"testing"
	"time"
)

func mockChecker(url string) bool {
	time.Sleep(50 * time.Millisecond)
	if url == "http://bad.omg" {
		return false
	}
	return true
}

func TestCheckWebsites(t *testing.T) {
	urls := []string{
		"http://good.com",
		"http://better.com",
		"http://bad.omg",
	}
	expected := map[string]bool{
		"http://good.com":   true,
		"http://better.com": true,
		"http://bad.omg":    false,
	}

	got := CheckWebsites(mockChecker, urls)

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func BenchmarkCheckWebsites(b *testing.B) {
	urls := make([]string, 100)
	for i := 0; i < len(urls); i++ {
		urls[i] = "a url"
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckWebsites(mockChecker, urls)
	}
}
