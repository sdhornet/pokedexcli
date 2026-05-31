package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sdhornet/pokedexcli/internal/pokecache"
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

func WalkMap(url string, cache *pokecache.Cache) (LocationData, error) {
	locations, err := getLocations(url, cache)
	if err != nil {
		return LocationData{}, err
	}
	return locations, nil
}

func getLocations(url string, cache *pokecache.Cache) (LocationData, error) {
	data, ok := cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return LocationData{}, fmt.Errorf("getting location data: %w", err)
		}
		defer res.Body.Close()

		if res.StatusCode < 200 || res.StatusCode > 299 {
			return LocationData{}, fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return LocationData{}, fmt.Errorf("reading response body: %w", err)
		}
		cache.Add(url, data)
	}

	var locations LocationData
	if err := json.Unmarshal(data, &locations); err != nil {
		return LocationData{}, fmt.Errorf("decoding location data: %w", err)
	}

	return locations, nil
}
