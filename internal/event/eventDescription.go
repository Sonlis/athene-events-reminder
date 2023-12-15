package event

import "fmt"

func (i *Ilmo) BuildEventDescription(eventId string) (string, error) {
	event, err := i.GetSingleEvent(eventId)
	if err != nil {
		return "", err
	}
	fmt.Println(event)
	return "Registration for " + event.Title + " starts in 5 minutes at " + i.WebURL + "/" + eventId, nil
}
