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

type LocationDetails struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name           string `json:"name"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	BaseExperience int    `json:"base_experience"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func WalkMap(url string, cache *pokecache.Cache) (LocationData, error) {
	var locations LocationData

	dataBytes, err := getPokeData(url, cache)
	if err != nil {
		return LocationData{}, err
	}

	if err := json.Unmarshal(dataBytes, &locations); err != nil {
		return LocationData{}, fmt.Errorf("decoding location data: %w", err)
	}

	return locations, nil
}

func ExploreLocation(url string, cache *pokecache.Cache) (LocationDetails, error) {
	var locationDetail LocationDetails

	dataBytes, err := getPokeData(url, cache)
	if err != nil {
		return LocationDetails{}, err
	}

	if err := json.Unmarshal(dataBytes, &locationDetail); err != nil {
		return LocationDetails{}, fmt.Errorf("decoding location details: %w", err)
	}

	return locationDetail, nil
}

func GatherPokemonDetails(url string, cache *pokecache.Cache) (Pokemon, error) {
	var pokemonDetails Pokemon

	dataBytes, err := getPokeData(url, cache)
	if err != nil {
		return Pokemon{}, err
	}

	if err := json.Unmarshal(dataBytes, &pokemonDetails); err != nil {
		return Pokemon{}, fmt.Errorf("decoding pokemon details: %w", err)
	}

	return pokemonDetails, nil
}

func getPokeData(url string, cache *pokecache.Cache) ([]byte, error) {
	data, ok := cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("getting PokeAPI data: %w", err)
		}
		defer res.Body.Close()

		if res.StatusCode < 200 || res.StatusCode > 299 {
			return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("reading response body: %w", err)
		}
		cache.Add(url, data)
	}

	return data, nil
}
