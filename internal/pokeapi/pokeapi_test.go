// ABOUTME: Integration tests for the pokeapi package against a local httptest server.
// ABOUTME: Covers the HTTP+decode round-trip, cache short-circuiting, and error paths — no real network.

package pokeapi

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sdhornet/pokedexcli/internal/pokecache"
)

const pikachuJSON = `{
	"name": "pikachu",
	"height": 4,
	"weight": 60,
	"base_experience": 112,
	"stats": [{"base_stat": 35, "stat": {"name": "hp"}}],
	"types": [{"type": {"name": "electric"}}]
}`

// GatherPokemonDetails should fetch over HTTP and decode the nested JSON into a
// Pokemon, including the stats/types slices whose struct tags are easy to get wrong.
func TestGatherPokemonDetailsDecodesResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(pikachuJSON))
	}))
	defer server.Close()

	cache := pokecache.NewCache(5 * time.Second)
	got, err := GatherPokemonDetails(server.URL, cache)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Name != "pikachu" {
		t.Errorf("Name = %q, want %q", got.Name, "pikachu")
	}
	if got.BaseExperience != 112 {
		t.Errorf("BaseExperience = %d, want %d", got.BaseExperience, 112)
	}
	if len(got.Stats) != 1 || got.Stats[0].BaseStat != 35 || got.Stats[0].Stat.Name != "hp" {
		t.Errorf("Stats = %+v, want one hp stat of 35", got.Stats)
	}
	if len(got.Types) != 1 || got.Types[0].Type.Name != "electric" {
		t.Errorf("Types = %+v, want one electric type", got.Types)
	}
}

// The second fetch of the same URL must be served from the cache, never touching
// the network again. This is the whole reason the pokecache layer exists.
func TestGatherPokemonDetailsCachesSecondCall(t *testing.T) {
	hits := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hits++
		w.Write([]byte(pikachuJSON))
	}))
	defer server.Close()

	cache := pokecache.NewCache(5 * time.Second)

	first, err := GatherPokemonDetails(server.URL, cache)
	if err != nil {
		t.Fatalf("first call errored: %v", err)
	}
	if hits != 1 {
		t.Fatalf("after first call, server hits = %d, want 1", hits)
	}

	second, err := GatherPokemonDetails(server.URL, cache)
	if err != nil {
		t.Fatalf("second call errored: %v", err)
	}
	if hits != 1 {
		t.Errorf("after second call, server hits = %d, want 1 (should be cached)", hits)
	}
	if !reflect.DeepEqual(first, second) {
		t.Errorf("cached result differs: first=%+v second=%+v", first, second)
	}
}

// A non-2xx response must surface as an error rather than a zero-value struct.
func TestGetPokeDataErrorsOnNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "missing", http.StatusNotFound)
	}))
	defer server.Close()

	cache := pokecache.NewCache(5 * time.Second)
	_, err := GatherPokemonDetails(server.URL, cache)
	if err == nil {
		t.Fatal("expected an error for a 404 response, got nil")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("error = %q, want it to mention the 404 status", err.Error())
	}
}

// A 200 response with an unparseable body must surface as a decode error.
func TestGatherPokemonDetailsErrorsOnBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("this is not json"))
	}))
	defer server.Close()

	cache := pokecache.NewCache(5 * time.Second)
	_, err := GatherPokemonDetails(server.URL, cache)
	if err == nil {
		t.Fatal("expected a decode error for a non-JSON body, got nil")
	}
}
