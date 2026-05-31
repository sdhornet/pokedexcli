## Assignment

Add an `explore` command. It takes the name of a location area as an argument.

**Run and submit** the CLI tests.

## Tips

- Use the same [PokeAPI location-area endpoint](https://pokeapi.co/docs/v2#location-areas), but this time you'll need to pass the `name` of the location area being explored. By adding a `name` or `id`, the API will return _a lot_ more information about the location area.
- Feel free to use tools like JSON lint and JSON to Go to help you parse the response.
- Parse the Pokemon's names from the response and display them to the user.
- Make sure to use the caching layer again! Re-exploring an area should be blazingly fast.
- You'll need to alter the function signature of _all_ your commands to allow them to allow parameters. E.g. `explore <area_name>`

**Example usage:**

```bash
Pokedex > explore pastoria-city-area
Exploring pastoria-city-area...
Found Pokemon:
 - tentacool
 - tentacruel
 - magikarp
 - gyarados
 - remoraid
 - octillery
 - wingull
 - pelipper
 - shellos
 - gastrodon
Pokedex >
```

## TODO
- Update the CLI to accept a second parameter
  - Struct?
- Update all the function signatures to accept parameters
- Most will just not use it, but for explore, validate it is real?
- Need a new struct of json data
- Integrate it with the get request and cacheing
  - Can probably reuse functions and pass the location to append to the url  
