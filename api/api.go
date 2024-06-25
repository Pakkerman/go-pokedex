package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pakkermandev/go-pokedex/pokecache"
)

type Location struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

var (
	Cache pokecache.Cache = *pokecache.NewCache(time.Millisecond * 5000)
	Page  int             = -1
)

func GetNextMap() {
	Page++
	endpoint := fmt.Sprintf("https://pokeapi.co/api/v2/location/?offset=%d", Page*20)

	cache, ok := Cache.Get(endpoint)
	if !ok {
		resp, err := http.Get(endpoint)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return
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
		return
	}

	var out strings.Builder
	for i := 0; i < len(location.Results); i++ {
		out.WriteString(location.Results[i].Name)
		out.WriteString("\n")
	}
	fmt.Printf("\n%vpage: %v\n", out.String(), Page+1)
}

func GetPreviousMap() {
	Page--
	endpoint := fmt.Sprintf("https://pokeapi.co/api/v2/location/?offset=%d", Page*20)

	cache, ok := Cache.Get(endpoint)
	if !ok {
		resp, err := http.Get(endpoint)
		if err != nil {
			fmt.Println("something wrong with GET\n", err)
			return
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return
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
		return
	}

	var out strings.Builder
	for i := 0; i < len(location.Results); i++ {
		out.WriteString(location.Results[i].Name)
		out.WriteString("\n")
	}
	fmt.Printf("\n%vpage: %v\n", out.String(), Page+1)
}

func Explore(name string) {
	endpoint := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", name)
	fmt.Printf("\nexploreing area %s\n", name)

	cache, ok := Cache.Get(endpoint)
	if !ok {
		resp, err := http.Get(endpoint)
		if err != nil {
			fmt.Println("something wrong with GET\n", err)
			return
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return
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
		return
	}

	var out strings.Builder
	for i := 0; i < len(location.PokemonEncounters); i++ {
		str := fmt.Sprintf("- %s\n", location.PokemonEncounters[i].Pokemon.Name)
		out.WriteString(str)
	}

	fmt.Println(out.String())
}

func PrintLocations(location []Location) {
	var out strings.Builder
	for i := 0; i < len(location); i++ {
		out.WriteString(location[i].Name)
		out.WriteString("\n")
	}
	fmt.Printf("\n%v\n", out.String())
}
