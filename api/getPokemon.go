package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
