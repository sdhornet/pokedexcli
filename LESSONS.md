# Lessons — Pokedex CLI

A running reference of higher-level Go techniques picked up in this project. Quick hits, not tutorials.

---

## Reading input line-by-line with `bufio.Scanner`

`os.Stdin` is an `io.Reader`. Wrapping it in a `bufio.Scanner` gives you a clean line-at-a-time loop.

```go
scanner := bufio.NewScanner(os.Stdin)
for {
    fmt.Print("pokedex > ")
    if !scanner.Scan() {
        break
    }
    line := scanner.Text()
    // ...
}
if err := scanner.Err(); err != nil {
    fmt.Fprintln(os.Stderr, "reading standard input:", err)
}
```

Key points:
- `Scan()` returns `false` on EOF **or** error. The two are indistinguishable inside the loop — check `Err()` after.
- `Text()` returns the line with the trailing newline stripped.
- Default split is by line. Swap with `scanner.Split(bufio.ScanWords)` to tokenize as it reads.
- Scanner has a max line length (64 KiB by default). For pasted blobs use `bufio.Reader` + `ReadString('\n')`.

## Normalizing user input

`strings.Fields` is what you almost always want over `strings.Split(s, " ")`:

```go
strings.Fields("  hello   world  ") // ["hello", "world"]
strings.Split ("  hello   world  ", " ") // ["", "", "hello", "", "", "world", "", ""]
```

`Fields` splits on any Unicode whitespace and collapses runs. Combine with `strings.ToLower` for case-insensitive command parsing.

## Map-of-structs as a command dispatch table

Instead of a giant `switch` for commands, store them in a map keyed by name. Each value is a struct that bundles metadata with the callback:

```go
type cliCommand struct {
    name        string
    description string
    callback    func() error
}

commands := map[string]cliCommand{
    "help": {name: "help", description: "...", callback: commandHelp},
    "exit": {name: "exit", description: "...", callback: commandExit},
}
```

Why this is nice:
- Adding a command is a single map entry — no edits to the dispatch logic.
- `help` can iterate the same map to print its own listing. The command list is data, not code.
- Functions are **first-class values** in Go — you can store them in structs, pass them around, return them from other functions.

The cost: ordering. Map iteration is randomized in Go, so `help` output shuffles between runs. If you want stable order, keep a separate `[]string` of keys or sort them in the iterator.

## The comma-ok idiom for map lookups

```go
elem, ok := commands[command[0]]
if !ok {
    fmt.Println("Unknown command")
    continue
}
```

A bare `commands[key]` returns the zero value when the key is missing — no error. The two-value form (`value, ok`) is how you distinguish "present and zero" from "absent." Same pattern shows up in type assertions and channel receives.

## Error handling style

Idiomatic Go:
- Functions return `error` as the **last** return value.
- Callers check `if err != nil` immediately.
- User-facing error messages go to `os.Stderr`, not stdout — so they don't get mixed with program output if someone pipes it.

```go
if err := elem.callback(); err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}
```

`%v` is the default formatter for any value, including errors. For wrapping, `fmt.Errorf("doing X: %w", err)` lets callers `errors.Is` / `errors.As` against the underlying error.

## `os.Exit` skips deferred functions

```go
func commandExit() error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil // unreachable
}
```

`os.Exit` terminates immediately — it does **not** run `defer`s. That's fine when there's nothing to clean up, but if you later open files, network connections, or anything else that needs flushing, prefer returning a sentinel from the callback and letting `main` exit normally. Pattern:

```go
type cliCommand struct {
    // ...
    callback func() error
}

// callback returns a sentinel error or a bool signaling exit,
// and main breaks the loop on it.
```

## Table-driven tests

The standard Go test pattern: a slice of anonymous structs, one per case, looped through.

```go
cases := []struct {
    input    string
    expected []string
}{
    {input: "  hello  world  ", expected: []string{"hello", "world"}},
    {input: "Charmander Bulbasaur PIKACHU", expected: []string{"charmander", "bulbasaur", "pikachu"}},
}

for _, c := range cases {
    actual := cleanInput(c.input)
    // assertions...
}
```

`t.Errorf` reports failure but lets the loop continue to the next case — good for seeing all failures at once. `t.Fatalf` halts the current test function. Inside a loop, `Errorf` + `continue` (as you have) is the usual move so one bad case doesn't mask others.

For larger suites, wrap each case in `t.Run(name, func(t *testing.T) { ... })` — gives each case its own subtest name in output and lets you target one with `go test -run TestCleanInput/case_name`.

## Splitting a small CLI into files

Three responsibilities, three places:

1. **Entrypoint** — `main.go` is a one-liner that calls into the REPL. No logic.
2. **REPL** — `repl.go` owns the loop, input parsing, dispatch, and the registry (`cliCommand` type + `getCommands()`).
3. **Commands** — one file per command (`command_exit.go`, `command_help.go`, ...). Each defines its callback. Same package as the registry, so no imports needed between them.

```
pokedexcli/
├── main.go            // calls startRepl()
├── repl.go            // loop + cliCommand + getCommands()
├── command_exit.go
└── command_help.go
```

Adding a new command = a new file + one line in `getCommands()`. Scales without a junk-drawer `commands.go`.

## Packages and the `internal/` directory

**A package = a directory.** All `.go` files in the directory declare the same `package <name>` at the top. Sibling files in the same package see each other with no imports — that's why `repl.go` can call `commandExit` directly.

**Exported = Capitalized.** Across package boundaries, only `Capitalized` identifiers are visible. `GetLocations` is exported; `getLocations` is package-private. Enforced by the compiler, not convention.

**`internal/` is special.** Packages under `internal/` can only be imported by code in the same module. Toolchain-enforced. Use it for implementation-detail packages you don't want imported from outside the project.

**Import path = module path + directory path.** Qualify calls with the package's declared name.

```go
// go.mod: module github.com/sdhornet/pokedexcli

// internal/pokeapi/pokeapi.go
package pokeapi
func GetLocations() { ... }
```

```go
// main.go
import "github.com/sdhornet/pokedexcli/internal/pokeapi"

pokeapi.GetLocations()
```

Directory name, `package` declaration, and the qualifier you call with are always the same string by convention. Lowercase, no underscores.

## `ABOUTME:` file headers (personal convention)

From global CLAUDE.md: every new source file starts with a 2-line comment block, both lines prefixed `ABOUTME:`, so a `grep "ABOUTME:"` across a project gives you a file-by-file map. Not a Go idiom — a Nate idiom.

---

*Add to this file as new techniques come up. Date entries if useful for spaced review.*
