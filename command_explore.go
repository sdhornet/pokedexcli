package main

import (
	"fmt"

	"github.com/sdhornet/pokedexcli/internal/pokeapi"
)

func explore(cfg *config, location string) error {
	baseURL := "https://pokeapi.co/api/v2/location-area/"

	if location == "" {
		return fmt.Errorf("No location provided to explore")
	}
	details, err := pokeapi.ExploreLocation(baseURL+location, cfg.Cache)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, v := range details.PokemonEncounters {
		fmt.Println(v.Pokemon.Name)
	}
	return nil
}
