package main

import (
	"fmt"

	"github.com/sdhornet/pokedexcli/internal/pokeapi"
)

func commandMap(cfg *config) error {
	if cfg.Next == "" {
		fmt.Println("your on the last page")
		return nil
	}
	locations, err := pokeapi.WalkMap(cfg.Next)
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
