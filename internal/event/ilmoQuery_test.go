package event

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	os.Setenv("ILMO_API_URL", "http://localhost:8080")
	i, err := Init()
	if err != nil {
		t.Errorf("Init() failed")
	}
	want := Ilmo{
		ApiURL: "http://localhost:8080",
	}
	if i != want {
		t.Errorf("Init() returned %+v, want %+v", i, want)
	}
}

func TestGetEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/events" {
			t.Errorf("Expected to request '/events', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
   {
      "id":1,
      "title":"test event 1",
      "date":"2023-12-12T14:45:00.000Z",
      "registrationStartDate":"2023-11-24T09:00:00.000Z",
      "quota":[
         {
            "title":"test quota 1",
            "size":25,
            "sortId":1,
            "signupCount":23
         },
         {
            "title":"test quota 2",
            "size":10,
            "sortId":2,
            "signupCount":7
         }
      ]
   },
   {
      "id":2,
      "title":"test event 2",
      "date":"2023-12-11T16:30:00.000Z",
      "registrationStartDate":"2023-12-08T16:00:00.000Z",
      "quota":[
         {
            "title":"test quota 1",
            "size":20,
            "sortId":1,
            "signupCount":0
         },
         {
            "title":"test quota 2",
            "size":20,
            "sortId":2,
            "signupCount":0
         },
         {
            "title":"test quota 3",
            "size":2,
            "sortId":3,
            "signupCount":0
         }
      ]
   },`))
	}))
	defer server.Close()
	os.Setenv("ILMO_API_URL", server.URL+"/events")
	i, err := Init()
	if err != nil {
		t.Errorf("Init() failed")
	}
	registration_date_1 := time.Date(2023, 11, 24, 9, 0, 0, 0, time.UTC)
	registration_date_2 := time.Date(2023, 12, 8, 16, 0, 0, 0, time.UTC)
	want := []Event{
		{
			ID:                    1,
			Title:                 "test event 1",
			RegistrationStartDate: registration_date_1,
		},
		{
			ID:                    2,
			Title:                 "test event 2",
			RegistrationStartDate: registration_date_2,
		},
	}
	got, err := i.GetEvents()
	if err != nil {
		t.Errorf("GetEvents() failed")
	}
	for i, event := range got {
		if event != want[i] {
			t.Errorf("GetEvents() returned %+v, want %+v", got, want)
		}
	}
}

func TestGetSingleEvent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/events/1" {
			t.Errorf("Expected to request '/events/1', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
   {
      "id":1,
      "title":"test event 1",
      "registrationStartDate":"2023-11-24T09:00:00.000Z",
      "quota":[
         {
            "title":"test quota 1",
            "size":25,
            "sortId":1,
            "signupCount":23
         },
         {
            "title":"test quota 2",
            "size":10,
            "sortId":2,
            "signupCount":7
         }
      ]
   }`))
	}))
	defer server.Close()
	os.Setenv("ILMO_API_URL", server.URL+"/events")
	i, err := Init()
	if err != nil {
		t.Errorf("Init() failed")

	}
	registration_date := time.Date(2023, 11, 24, 9, 0, 0, 0, time.UTC)
	want := Event{
		ID:                    1,
		Title:                 "test event 1",
		RegistrationStartDate: registration_date,
	}
	got, err := i.GetSingleEvent("1")
	if err != nil {
		t.Errorf("getSingleEvent() failed")
	}
	if got != want {
		t.Errorf("getSingleEvent() returned %+v, want %+v", got, want)
	}
}
