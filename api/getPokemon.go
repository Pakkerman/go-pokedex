package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pakkermandev/go-pokedex/pokecache"
)

func GetPokemon(name string) (Pokemon, error) {
	var pokemon Pokemon

	cache, ok := pokecache.PokeCache.Get(name)
	if ok {
		if err := json.Unmarshal(cache, &pokemon); err != nil {
			return pokemon, nil
		}

		return pokemon, nil
	}

	endpoint := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)
	resp, err := http.Get(endpoint)
	if err != nil {
		return pokemon, nil
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return pokemon, nil
	}

	pokecache.PokeCache.Add(name, body)

	if err := json.Unmarshal(body, &pokemon); err != nil {
		return pokemon, err
	}

	return pokemon, nil
}
