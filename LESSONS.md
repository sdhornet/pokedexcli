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

Practical rule:
- Use `return err` when the caller already has enough context.
- Use `fmt.Errorf("doing thing: %w", err)` when adding context makes the failure easier to locate.
- Use `fmt.Errorf("bad status code: %d", status)` for errors you create yourself with no underlying error to wrap.
- Prefer lowercase error strings. The caller may print a prefix like `Error:`, so `fmt.Errorf("getting location data: %w", err)` reads better than `fmt.Errorf("Error getting location data: %w", err)`.

```go
res, err := http.Get(url)
if err != nil {
    return LocationData{}, fmt.Errorf("getting location data: %w", err)
}
defer res.Body.Close()

if res.StatusCode < 200 || res.StatusCode > 299 {
    return LocationData{}, fmt.Errorf("unexpected status code: %d", res.StatusCode)
}
```

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

## Shared REPL state with a config pointer

The REPL needs memory between commands. For example, `map` gets a `next` URL from the API response, and `mapb` needs to remember the `previous` URL.

Store that shared state in one config value created when the REPL starts:

```go
type config struct {
    Next     string
    Previous string
}

func startRepl() {
    cfg := &config{Next: "https://pokeapi.co/api/v2/location-area/"}
    // ...
}
```

Then make every command callback accept the same pointer:

```go
type cliCommand struct {
    name        string
    description string
    callback    func(*config) error
}
```

Why a pointer? Without one, each command would receive a copy of the config. Updating the copy would not update the REPL's shared state. With `*config`, commands can update the original:

```go
cfg.Next = locations.Next
cfg.Previous = locations.Previous
```

Because `cfg` is a pointer to a struct, Go automatically dereferences it for field access. This:

```go
cfg.Next = locations.Next
```

is shorthand for:

```go
(*cfg).Next = locations.Next
```

Same behavior applies when reading fields.

## Pagination state belongs at the command/REPL layer

PokeAPI pagination responses include `next` and `previous` URLs. The API package can return those values, but the REPL config owns the current navigation state.

Flow:

```text
map command
  -> use cfg.Next
  -> fetch location data
  -> print location names
  -> update cfg.Next and cfg.Previous

mapb command
  -> if cfg.Previous is empty, tell the user they're on the first page
  -> otherwise use cfg.Previous
  -> fetch location data
  -> print location names
  -> update cfg.Next and cfg.Previous
```

Seeding the first URL in `startRepl` is cleaner than making `commandMap` guess whether an empty URL means "first run" or "no next page".

## HTTP request and JSON decoding flow

The PokeAPI package owns the network and JSON details. The command layer should not need to know how to make an HTTP request or decode a response body.

Basic shape:

```go
res, err := http.Get(url)
if err != nil {
    return LocationData{}, fmt.Errorf("getting location data: %w", err)
}
defer res.Body.Close()

if res.StatusCode < 200 || res.StatusCode > 299 {
    return LocationData{}, fmt.Errorf("unexpected status code: %d", res.StatusCode)
}

data, err := io.ReadAll(res.Body)
if err != nil {
    return LocationData{}, fmt.Errorf("reading response body: %w", err)
}

var locations LocationData
if err := json.Unmarshal(data, &locations); err != nil {
    return LocationData{}, fmt.Errorf("decoding location data: %w", err)
}
```

Key points:
- `defer res.Body.Close()` should happen after confirming `http.Get` returned no error, before any early returns based on status code.
- `res.Body` is a stream, not a string. To inspect the body text, read it with `io.ReadAll`.
- Check status codes before assuming the body is valid JSON.
- Use struct tags like `` `json:"next"` `` to map JSON fields onto Go struct fields.

## Exported response types across package boundaries

If the `main` package needs to use fields from a type returned by `pokeapi`, the type and fields need to be exported.

```go
type LocationData struct {
    Count    int    `json:"count"`
    Next     string `json:"next"`
    Previous string `json:"previous"`
    Results  []struct {
        Name string `json:"name"`
        URL  string `json:"url"`
    } `json:"results"`
}
```

`LocationData`, `Next`, `Previous`, and `Results` are capitalized, so code outside the `pokeapi` package can use them. A function can technically return an unexported type, but that is awkward for callers because they cannot name the type in their own signatures.

## Caching raw API responses

The cache package stores raw response bytes by URL:

```text
URL string -> []byte response body
```

It does not know anything about Pokemon, locations, JSON structs, or the CLI. That keeps `internal/pokecache` reusable. The PokeAPI package decides what the bytes mean by unmarshalling them into `LocationData`.

```go
type cacheEntry struct {
    createdAt time.Time
    val       []byte
}

type Cache struct {
    entries map[string]cacheEntry
    mu      sync.Mutex
}
```

`NewCache` initializes the map and starts cleanup in the background:

```go
func NewCache(interval time.Duration) *Cache {
    cache := &Cache{entries: make(map[string]cacheEntry)}
    go cache.reapLoop(interval)
    return cache
}
```

The map must be made with `make`; writing to a nil map panics.

## Pointer receivers and mutexes

Cache methods use pointer receivers:

```go
func (c *Cache) Add(key string, val []byte) { ... }
func (c *Cache) Get(key string) ([]byte, bool) { ... }
```

This matters because `Cache` contains a `sync.Mutex`. Copying structs that contain mutexes is a bad habit: different copies can appear to protect the same data but actually use different locks. Pointer receivers keep every method operating on the same cache and same mutex.

The basic lock pattern:

```go
c.mu.Lock()
defer c.mu.Unlock()
```

Use this around every map access: adding, getting, and deleting. Go maps are not safe for concurrent read/write access.

## Adding and getting cache entries

For `Add`, store a full `cacheEntry` in the map:

```go
c.entries[key] = cacheEntry{
    createdAt: time.Now(),
    val:       val,
}
```

For `Get`, use the comma-ok idiom:

```go
entry, ok := c.entries[key]
if !ok {
    return nil, false
}

return entry.val, true
```

Returning `nil, false` is normal for a missing `[]byte`; `nil` is the zero value for slices.

Map values are not directly addressable. If you need to modify part of a struct already stored in a map, pull the struct out, modify the copy, then store it back:

```go
entry := c.entries[key]
entry.val = newVal
c.entries[key] = entry
```

For `Add`, replacing the whole entry with a struct literal is simpler.

## Background cleanup with tickers and goroutines

`time.NewTicker(interval)` creates a value that sends a signal on `ticker.C` every interval.

```go
func (c *Cache) reapLoop(interval time.Duration) {
    ticker := time.NewTicker(interval)
    for range ticker.C {
        c.mu.Lock()

        for key, entry := range c.entries {
            if time.Since(entry.createdAt) > interval {
                delete(c.entries, key)
            }
        }

        c.mu.Unlock()
    }
}
```

`go cache.reapLoop(interval)` starts that loop in the background. `NewCache` returns immediately while the cleanup loop keeps running.

`delete(c.entries, key)` is the built-in way to remove a map entry.

## Integrating cache into the request layer

The REPL config owns one cache for the whole session:

```go
type config struct {
    Next     string
    Previous string
    Cache    *pokecache.Cache
}
```

Create it once:

```go
cfg := &config{
    Next:  "https://pokeapi.co/api/v2/location-area/",
    Cache: pokecache.NewCache(5 * time.Second),
}
```

Commands pass `cfg.Cache` into the PokeAPI request layer. The request layer checks the cache before making an HTTP request:

```text
getLocations(url, cache)
  -> cache.Get(url)
     -> if hit: use cached bytes
     -> if miss: http.Get(url), read body, cache.Add(url, bytes)
  -> json.Unmarshal(bytes, &locations)
```

Cached bytes and fresh HTTP bytes both go through the same JSON decoding path.

## Testing cache behavior with fake data

Cache tests should use fake keys and byte slices, not real PokeAPI responses. The cache package's responsibility is only `string -> []byte` storage.

```go
cache := NewCache(5 * time.Second)
cache.Add("https://example.com", []byte("testdata"))

val, ok := cache.Get("https://example.com")
```

The add/get test proves values can be stored and retrieved.

The reap-loop test uses a tiny interval and `time.Sleep`:

```go
cache := NewCache(5 * time.Millisecond)
cache.Add("https://example.com", []byte("testdata"))

time.Sleep(10 * time.Millisecond)

_, ok := cache.Get("https://example.com")
```

This keeps the test fast while proving old entries are removed.

## `ABOUTME:` file headers (personal convention)

From global CLAUDE.md: every new source file starts with a 2-line comment block, both lines prefixed `ABOUTME:`, so a `grep "ABOUTME:"` across a project gives you a file-by-file map. Not a Go idiom — a Nate idiom.

---

*Add to this file as new techniques come up. Date entries if useful for spaced review.*
