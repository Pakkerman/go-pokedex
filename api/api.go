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

type Pokemon struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		}
	} `json:"types"`
}

var (
	Cache pokecache.Cache = *pokecache.NewCache(time.Millisecond * 5000)
	Page  int             = -1
)

func GetNextMap() error {
	Page++
	endpoint := fmt.Sprintf("https://pokeapi.co/api/v2/location/?offset=%d", Page*20)

	cache, ok := Cache.Get(endpoint)
	if !ok {
		resp, err := http.Get(endpoint)
		if err != nil {
			fmt.Println(err)
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
		Results []struct {
			Name string `json:"name"`
		} `json:"results"`
	}{}

	if err := json.Unmarshal(cache, &location); err != nil {
		return err
	}

	var out strings.Builder
	for i := 0; i < len(location.Results); i++ {
		out.WriteString(location.Results[i].Name)
		out.WriteString("\n")
	}
	fmt.Printf("\n%vpage: %v\n", out.String(), Page+1)
	return nil
}

func GetPreviousMap() error {
	Page--
	endpoint := fmt.Sprintf("https://pokeapi.co/api/v2/location/?offset=%d", Page*20)

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
		Results []struct {
			Name string `json:"name"`
		} `json:"results"`
	}{}

	if err := json.Unmarshal(cache, &location); err != nil {
		return err
	}

	var out strings.Builder
	for i := 0; i < len(location.Results); i++ {
		out.WriteString(location.Results[i].Name)
		out.WriteString("\n")
	}
	fmt.Printf("\n%vpage: %v\n", out.String(), Page+1)
	return nil
}

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

// func Catch(name string) (Pokemon, error) {
// 	var pokemon Pokemon
//
// 	cacheEntry, ok := Cache.Get(name)
// 	if ok {
//
// 		if err := json.Unmarshal(cacheEntry, &pokemon); err != nil {
// 			return pokemon, err
// 		}
// 		return pokemon, nil
// 	}
//
// 	endpoint := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)
// 	resp, err := http.Get(endpoint)
// 	if err != nil {
// 		return pokemon, err
// 	}
//
// 	defer resp.Body.Close()
//
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return pokemon, err
// 	}
//
// 	Cache.Add(name, body)
//
// 	if err := json.Unmarshal(body, &pokemon); err != nil {
// 		return pokemon, err
// 	}
//
// 	return pokemon, nil
// }

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

func GetPokemon(name string) (Pokemon, error) {
	var pokemon Pokemon

	cache, ok := Cache.Get(name)
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

	Cache.Add(name, body)

	if err := json.Unmarshal(body, &pokemon); err != nil {
		return pokemon, err
	}

	return pokemon, nil
}
