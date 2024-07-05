package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pakkermandev/go-pokedex/pokecache"
)

func GetLocations(page int) (Locations, error) {
	var locations Locations

	endpoint := fmt.Sprintf("https://pokeapi.co/api/v2/location/?offset=%d", page*20)
	cache, ok := pokecache.PokeCache.Get(endpoint)
	if ok {
		if err := json.Unmarshal(cache, &locations); err != nil {
			return locations, err
		}
		return locations, nil
	}

	resp, err := http.Get(endpoint)
	if err != nil {
		return locations, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return locations, err
	}

	pokecache.PokeCache.Add(endpoint, body)

	if err := json.Unmarshal(body, &locations); err != nil {
		return locations, err
	}

	return locations, nil
}
