package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocationData struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func WalkMap(url string) (LocationData, error) {
	locations, err := getLocations(url)
	if err != nil {
		return LocationData{}, err
	}
	return locations, nil
}

func getLocations(url string) (LocationData, error) {
	res, err := http.Get(url)
	if err != nil {
		return LocationData{}, fmt.Errorf("getting location data: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return LocationData{}, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationData{}, fmt.Errorf("reading response body: %w", err)
	}

	var locations LocationData
	if err := json.Unmarshal(data, &locations); err != nil {
		return LocationData{}, fmt.Errorf("decoding location data: %w", err)
	}

	return locations, nil
}
