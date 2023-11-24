package utils

import (
	"io"
	"net/http"
)

// FetchURL fetches the contents of a URL and returns them as a byte array.
func FetchURL(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
