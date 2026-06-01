package main

import (
	"fmt"
	"math/rand"

	"github.com/sdhornet/pokedexcli/internal/pokeapi"
)

func calculateCatchChance(baseExperince int) bool {
	catchChance := 100 - baseExperince/2

	if catchChance < 10 {
		catchChance = 10
	}
	if catchChance > 90 {
		catchChance = 90
	}

	roll := rand.Intn(100)

	return roll < catchChance
}

func commandCatch(cfg *config, pokemon string) error {
	baseURL := "https://pokeapi.co/api/v2/pokemon/"

	if pokemon == "" {
		return fmt.Errorf("no pokemon provided to catch")
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)

	pokemonDetails, err := pokeapi.GatherPokemonDetails(baseURL+pokemon, cfg.Cache)
	if err != nil {
		return err
	}

	catch := calculateCatchChance(pokemonDetails.BaseExperience)

	if !catch {
		fmt.Printf("%s escaped!\n", pokemonDetails.Name)
		return nil
	}

	cfg.Pokedex[pokemonDetails.Name] = pokemonDetails
	fmt.Printf("%s was caught!\n", pokemonDetails.Name)

	return nil
}
