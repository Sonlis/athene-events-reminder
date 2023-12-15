package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchURL(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			t.Errorf("Expected to request '/', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`oui`))
	}))
	defer server.Close()

	response, err := FetchURL(server.URL)
	if err != nil {
		t.Errorf("Error fetching test URL: %v", err)
	}

	want := []byte(`oui`)
	for position, byte := range response {
		if byte != want[position] {
			t.Errorf("Fetching the test URL returned %s, want %v", response, want)
		}
	}
}
