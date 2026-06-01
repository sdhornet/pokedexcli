## Assignment

- [ ] **Add a `catch` command**. It takes the name of a Pokemon as an argument. _Example usage_:

```bash
Pokedex > catch pikachu
Throwing a Pokeball at pikachu...
pikachu escaped!
Pokedex > catch pikachu
Throwing a Pokeball at pikachu...
pikachu was caught!
```

- [ ] Be sure to print the `Throwing a Pokeball at <pokemon>...` message before determining if the Pokemon was caught or not.
- [ ] Use the [Pokemon endpoint](https://pokeapi.co/docs/v2#pokemon) to get information about a Pokemon by name.
- [ ] Give the user a _chance_ to catch the Pokemon using the [math/rand package](https://pkg.go.dev/math/rand#Rand.Intn).
- [ ] You can use the pokemon's "base experience" to determine the chance of catching it. The higher the base experience, the harder it should be to catch.
- [ ] Once the Pokemon is caught, add it to the user's Pokedex. I used a `map[string]Pokemon` to keep track of caught Pokemon.
- [ ] Test the `catch` command manually - make sure you can actually catch a Pokemon within a reasonable number of tries.

**Run and submit** the CLI tests.