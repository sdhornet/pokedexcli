package main

import (
	"fmt"

	"github.com/sdhornet/pokedexcli/internal/pokeapi"
)

func commandMapb(cfg *config, location string) error {
	if cfg.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	locations, err := pokeapi.WalkMap(cfg.Previous, cfg.Cache)
	if err != nil {
		return err
	}
	cfg.Previous = locations.Previous
	cfg.Next = locations.Next

	for _, v := range locations.Results {
		fmt.Println(v.Name)
	}
	return nil
}
