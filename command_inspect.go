package main

import (
	"fmt"
)

func commandInspect(cfg *config, pokemon string) error {
	if pokemon == "" {
		return fmt.Errorf("no pokemon provided to inspect")
	}

	elem, ok := cfg.Pokedex[pokemon]
	if !ok {
		fmt.Println("You have not caught that pokemon")
		return nil
	}

	fmt.Println("Name: ", elem.Name)
	fmt.Println("Height: ", elem.Height)
	fmt.Println("Weight: ", elem.Weight)
	fmt.Println("Stats:")
	for _, v := range elem.Stats {
		fmt.Printf("  - %s: %d\n", v.Stat.Name, v.BaseStat)
	}
	fmt.Println("Types:")
	for _, j := range elem.Types {
		fmt.Printf("  - %s\n", j.Type.Name)
	}

	return nil
}
