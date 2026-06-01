# Pokedex CLI

A small interactive command-line Pokedex written in Go as part of the
[Boot.dev](https://www.boot.dev/) backend curriculum.

The program uses the public [PokeAPI](https://pokeapi.co/) location-area
endpoints to page through Pokemon location areas and explore which Pokemon can
be encountered in a specific area. API responses are cached in memory so repeat
requests are fast during a REPL session.

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
