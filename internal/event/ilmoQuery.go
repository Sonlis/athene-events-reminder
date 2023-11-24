package event

import (
	"context"
	"encoding/json"
	"github.com/Sonlis/athene-events-notifier/internal/utils"
	"github.com/sethvargo/go-envconfig"
	"time"
)

type Ilmo struct {
	ApiURL string `env:"ILMO_API_URL"`
}

type Event struct {
	ID                    int       `json:"id"`
	Title                 string    `json:"title"`
	RegistrationStartDate time.Time `json:"registrationStartDate"`
}

func Init() (Ilmo, error) {
	ctx := context.Background()
	var i Ilmo
	err := envconfig.Process(ctx, &i)
	return i, err
}

func (i *Ilmo) GetEvents() ([]Event, error) {
	var events []Event
	eventsJson, err := utils.FetchURL(i.ApiURL)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(eventsJson), &events)
	return events, nil
}

func (i *Ilmo) GetSingleEvent(eventId string) (Event, error) {
	var event Event
	eventJson, err := utils.FetchURL(i.ApiURL + "/" + eventId)
	if err != nil {
		return event, err
	}
	json.Unmarshal([]byte(eventJson), &event)
	return event, nil
}
