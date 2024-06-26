package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/pakkermandev/go-pokedex/pokecache"
)

var Cache pokecache.Cache = *pokecache.NewCache(time.Millisecond * 5000)

func Explore(name string) error {
	var endpoint string
	match, err := regexp.MatchString("[0-9]+", name)
	if err != nil {
		return err
	}

	if match {
		endpoint = fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", name)
	} else {
		endpoint = fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s-area", name)
	}

	fmt.Printf("\nexploreing area %s\n", name)

	cache, ok := Cache.Get(endpoint)
	if !ok {
		resp, err := http.Get(endpoint)
		if err != nil {
			fmt.Println("something wrong with GET\n", err)
			return err
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		Cache.Add(endpoint, body)
		cache = body
	}

	location := struct {
		PokemonEncounters []struct {
			Pokemon struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"pokemon"`
			VersionDetails []struct {
				EncounterDetails []struct {
					Chance          int   `json:"chance"`
					ConditionValues []any `json:"condition_values"`
					MaxLevel        int   `json:"max_level"`
					Method          struct {
						Name string `json:"name"`
						URL  string `json:"url"`
					} `json:"method"`
					MinLevel int `json:"min_level"`
				} `json:"encounter_details"`
				MaxChance int `json:"max_chance"`
				Version   struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"version"`
			} `json:"version_details"`
		} `json:"pokemon_encounters"`
	}{}

	if err := json.Unmarshal(cache, &location); err != nil {
		return err
	}

	var out strings.Builder
	for i := 0; i < len(location.PokemonEncounters); i++ {
		str := fmt.Sprintf("- %s\n", location.PokemonEncounters[i].Pokemon.Name)
		out.WriteString(str)
	}

	fmt.Println(out.String())
	return nil
}

func Inspect(name string) (Pokemon, error) {
	var pokemon Pokemon

	cacheEntry, ok := Cache.Get(name)
	if !ok {
		endpoint := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)

		resp, err := http.Get(endpoint)
		if err != nil {
			return pokemon, err
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return pokemon, err
		}

		cacheEntry = body
	}

	if err := json.Unmarshal(cacheEntry, &pokemon); err != nil {
		return pokemon, err
	}

	return pokemon, nil
}
