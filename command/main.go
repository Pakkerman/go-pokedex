package command

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/pakkermandev/go-pokedex/api"
)

type CliCommand struct {
	Name        string
	Description string
	Callback    func(*string) error
}

type Pokedex struct {
	PokemonsCaught map[string]string
}

var (
	pokedex Pokedex = Pokedex{PokemonsCaught: make(map[string]string)}
	page    int     = -1
)

func GetOptions() map[string]CliCommand {
	return map[string]CliCommand{
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    commandHelp,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"map": {
			Name:        "map",
			Description: "Get the next 20 locations",
			Callback:    mapLocations,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Get the previous 20 locations",
			Callback:    mapbLocations,
		},
		"explore": {
			Name:        "explore",
			Description: "explore the area",
			Callback:    explore,
		},
		"catch": {
			Name:        "catch",
			Description: "catch a pokemon",
			Callback:    catch,
		},
		"inspect": {
			Name:        "catch",
			Description: "catch a pokemon",
			Callback:    inspect,
		},
	}
}

func commandHelp(arg *string) error {
	fmt.Println("\nUsage:")

	options := GetOptions()

	for _, opt := range options {
		fmt.Printf("%v:\t%v\n", opt.Name, opt.Description)
	}

	fmt.Println()
	return nil
}

func commandExit(arg *string) error {
	fmt.Println("Exit the program")
	os.Exit(0)
	return nil
}

func mapLocations(arg *string) error {
	page++
	locations, err := api.GetLocations(page)
	if err != nil {
		return err
	}

	var out strings.Builder
	for i := 0; i < len(locations.Results); i++ {
		out.WriteString(locations.Results[i].Name)
		out.WriteString("\n")
	}
	fmt.Printf("\n%vpage: %v\n", out.String(), page+1)

	return nil
}

func mapbLocations(arg *string) error {
	if page-1 <= 0 {
		fmt.Println("You are at the first page")
	} else {
		page--
	}

	locations, err := api.GetLocations(page)
	if err != nil {
		return err
	}

	var out strings.Builder
	for i := 0; i < len(locations.Results); i++ {
		out.WriteString(locations.Results[i].Name)
		out.WriteString("\n")
	}
	fmt.Printf("\n%vpage: %v\n", out.String(), page+1)
	return nil
}

func explore(arg *string) error {
	name := *arg

	err := api.Explore(name)
	if err != nil {
		return err
	}
	return nil
}

func catch(arg *string) error {
	name := *arg
	pokemon, err := api.GetPokemon(name)
	if err != nil {
		return err
	}

	randomNumber := rand.Intn(100)
	chance := pokemon.BaseExperience / 5
	if randomNumber < chance {
		fmt.Println("You fail to catch", pokemon.Name)
		return nil
	}

	fmt.Println("You successfully captured", pokemon.Name)
	pokedex.PokemonsCaught[pokemon.Name] = pokemon.Name

	return nil
}

func inspect(arg *string) error {
	name := *arg

	pokemonCaught, ok := pokedex.PokemonsCaught[name]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	pokemon, err := api.GetPokemon(pokemonCaught)
	if err != nil {
		return err
	}

	var out strings.Builder
	out.WriteString(fmt.Sprintf("Name: %s\n", pokemon.Name))
	out.WriteString(fmt.Sprintf("Height: %d\n", pokemon.Height))
	out.WriteString(fmt.Sprintf("Weight: %d\n", pokemon.Weight))

	out.WriteString("Stats:\n")
	for i := 0; i < len(pokemon.Stats); i++ {
		out.WriteString(fmt.Sprintf("  -%s: %d\n", pokemon.Stats[i].Stat.Name, pokemon.Stats[i].BaseStat))
	}

	out.WriteString("Types:\n")
	for i := 0; i < len(pokemon.Types); i++ {
		out.WriteString(fmt.Sprintf("  - %s", pokemon.Types[i].Type.Name))
	}

	out.WriteString("\n")
	fmt.Print(out.String())
	return nil
}
