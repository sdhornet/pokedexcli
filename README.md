# Pokedex CLI

A small interactive command-line Pokedex written in Go as part of the
[Boot.dev](https://www.boot.dev/) backend curriculum.

The program uses the public [PokeAPI](https://pokeapi.co/) endpoints to page
through Pokemon location areas, explore which Pokemon can be encountered in a
specific area, catch them, and inspect the ones you have added to your Pokedex.
API responses are cached in memory so repeat requests are fast during a REPL
session.

## Requirements

- Go, using the version declared in `go.mod`
- Network access to `https://pokeapi.co`

## Run

Start the REPL:

```bash
go run .
```

You should see a prompt:

```text
pokedex >
```

## Commands

### `help`

Prints the available commands.

```text
pokedex > help
```

### `map`

Displays the next page of location areas from PokeAPI.

```text
pokedex > map
canalave-city-area
eterna-city-area
pastoria-city-area
...
```

### `mapb`

Displays the previous page of location areas.

```text
pokedex > mapb
```

If you are already on the first page, the CLI prints a message instead of
making another request.

### `explore <area_name>`

Displays Pokemon encounters for a specific location area.

```text
pokedex > explore pastoria-city-area
Found Pokemon:
tentacool
tentacruel
magikarp
gyarados
remoraid
octillery
wingull
pelipper
shellos
gastrodon
```

Location area names come from the `map` and `mapb` commands.

### `catch <pokemon_name>`

Attempts to catch a Pokemon. The chance of success scales with the Pokemon's
base experience, so stronger Pokemon are harder to catch. Caught Pokemon are
stored in your Pokedex for the rest of the session.

```text
pokedex > catch pikachu
Throwing a Pokeball at pikachu...
pikachu was caught!
You may now inspect it with the inspect command.
```

A Pokemon that gets away prints `pikachu escaped!` instead.

### `inspect <pokemon_name>`

Shows the name, height, weight, stats, and types of a Pokemon you have already
caught.

```text
pokedex > inspect pikachu
Name:  pikachu
Height:  4
Weight:  60
Stats:
  - hp: 35
  - attack: 55
  - defense: 40
  - special-attack: 50
  - special-defense: 50
  - speed: 90
Types:
  - electric
```

If you have not caught that Pokemon yet, the CLI prints
`You have not caught that pokemon`.

### `pokedex`

Lists every Pokemon you have caught this session.

```text
pokedex > pokedex
Your Pokedex:
 - pikachu
 - charmander
```

Before you have caught anything it prints `You have not caught any pokemon yet!`.

### `exit`

Exits the REPL.

```text
pokedex > exit
Closing the Pokedex... Goodbye!
```

## Test

Run the project tests:

```bash
go test ./...
```

## Project Structure

- `main.go` starts the REPL.
- `repl.go` handles input parsing, command dispatch, and shared REPL state.
- `command_*.go` files implement individual commands.
- `internal/pokeapi` fetches and decodes PokeAPI responses.
- `internal/pokecache` provides the in-memory response cache.
