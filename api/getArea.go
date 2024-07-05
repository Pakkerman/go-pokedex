package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/pakkermandev/go-pokedex/pokecache"
)

func GetArea(name string) (Area, error) {
	var area Area
	match, err := regexp.MatchString("[0-9]+", name)
	if err != nil {
		return area, err
	}

	var endpoint string
	if match {
		endpoint = fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", name)
	} else {
		endpoint = fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s-area", name)
	}

	cache, ok := pokecache.PokeCache.Get(endpoint)
	if ok {
		if err := json.Unmarshal(cache, &area); err != nil {
			return area, err
		}
	}

	resp, err := http.Get(endpoint)
	if err != nil {
		return area, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return area, err
	}

	err = json.Unmarshal(body, &area)
	if err != nil {
		return area, err
	}

	return area, nil
}
