package main

import (
	"fmt"
)

func commandPokedex(cfg *config, _ string) error {
	if len(cfg.Pokedex) == 0 {
		fmt.Println("You have not caught any pokemon yet!")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for name, _ := range cfg.Pokedex {
		fmt.Printf(" - %s\n", name)
	}

	return nil
}
