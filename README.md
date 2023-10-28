# shakemon
Pokemon's description in Shakespeare style.

# About
This project has been a testing playground for:
  - evaluating [Codeium](https://codeium.com/) capabilities for Golang
  - getting my feet wet at using [resty](https://github.com/go-resty/resty) library as an http client.

# Usage
This is a Go project, installing a golang compiler is a must.

- clone the project by `go build`
- run `./shakemon`
- from a terminal: `curl localhost:8080/translate/pokemon/pikachu | jq` or any other celebrity Pokemon
- output by the Bard itself:
```json
{
  "name": "pikachu",
  "description": "At which hour several of these pokémon gather,  their electricity couldst buildeth and cause lightning storms."
}
```

(Note: `jq` is but a pretty printing of the json output; skip it if unpleasant)

# So... how's thnis AI engine Codeium faring?
Pretty darn good, I have been pleasantly impressed. 

Read the full [experience](https://lucawolf.github.io/shakemon/) 
of this mini experiment to see what it has been capable of.
