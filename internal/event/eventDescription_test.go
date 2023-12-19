package event

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"os"
)

func TestBuildEventDescription(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/1" {
			t.Errorf("Expected to request '/1', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
        {
            "id": 1,
            "title": "Test event"
        }
        `))
	}))
	defer server.Close()
	os.Setenv("ILMO_API_URL", server.URL)
	os.Setenv("ILMO_WEB_URL", server.URL)
	i, err := Init()
	if err != nil {
		t.Errorf("Failed to initialize test ilmo config: %v", err)
	}
	eventDescription, err := i.BuildEventDescription("1")
	if err != nil {
		t.Errorf("Failed to build event description: %v", err)
	}
	want := "Registration for Test event starts in 5 minutes at " + i.WebURL + "/1"
	if eventDescription != want {
		t.Errorf("Expected event description to be %s, got: %s", want, eventDescription)
	}
}
