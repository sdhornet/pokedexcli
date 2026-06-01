package main

import (
	"fmt"

	"github.com/sdhornet/pokedexcli/internal/pokeapi"
)

func commandMap(cfg *config, _ string) error {
	if cfg.Next == "" {
		fmt.Println("you're on the last page")
		return nil
	}
	locations, err := pokeapi.WalkMap(cfg.Next, cfg.Cache)
	if err != nil {
		return err
	}
	cfg.Next = locations.Next
	cfg.Previous = locations.Previous

	for _, v := range locations.Results {
		fmt.Println(v.Name)
	}
	return nil
}
