package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sdhornet/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	Next     string
	Previous string
	Cache    *pokecache.Cache
}

func startRepl() {
	cfg := &config{
		Next: "https://pokeapi.co/api/v2/location-area/",
		Cache: pokecache.NewCache(5 * time.Second),
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("pokedex > ")
		if !scanner.Scan() {
			break
		}

		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		command, ok := getCommands()[commandName]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		if err := command.callback(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations",
			callback:    commandMapb,
		},
	}
}
