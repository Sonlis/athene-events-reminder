package event

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuildEventDescription(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			t.Errorf("Expected to request '/', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`oui`))
	}))
	defer server.Close()
}
