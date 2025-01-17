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
			Description: "Map out the next 20 locations",
			Callback:    mapLocations,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Map out the previous 20 locations",
			Callback:    mapbLocations,
		},
		"explore": {
			Name:        "explore",
			Description: "Explore the area",
			Callback:    explore,
		},
		"catch": {
			Name:        "catch",
			Description: "Try to catch a pokemon",
			Callback:    catch,
		},
		"inspect": {
			Name:        "insepct",
			Description: "Inspect a pokemon that you have captured",
			Callback:    inspect,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "List all your pokemons",
			Callback:    getPokedex,
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
		str := fmt.Sprintf("%02d: %s\n", i+1+(page*20), locations.Results[i].Name)
		out.WriteString(str)
	}
	fmt.Printf("\n%vpage: %v\n", out.String(), page+1)

	return nil
}

func mapbLocations(arg *string) error {
	if page-1 < 0 {
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
		str := fmt.Sprintf("%02d: %s\n", i+1+(page*20), locations.Results[i].Name)
		out.WriteString(str)
	}
	fmt.Printf("\n%vpage: %v\n", out.String(), page+1)
	return nil
}

func explore(arg *string) error {
	name := *arg

	area, err := api.GetArea(name)
	if err != nil {
		return err
	}

	var out strings.Builder
	out.WriteString(fmt.Sprintf("\nExploring area: %s\n", area.Names[0].Name))
	for i := 0; i < len(area.PokemonEncounters); i++ {
		str := fmt.Sprintf("- %s\n", area.PokemonEncounters[i].Pokemon.Name)
		out.WriteString(str)
	}

	fmt.Println(out.String())
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
		out.WriteString(fmt.Sprintf("  - %s\n", pokemon.Types[i].Type.Name))
	}

	out.WriteString("\n")
	fmt.Print(out.String())
	return nil
}

func getPokedex(arg *string) error {
	if len(pokedex.PokemonsCaught) == 0 {
		fmt.Printf("You have not catch any pokemons!\n")
		return nil

	}

	var out strings.Builder

	out.WriteString("You have caught: \n")
	for key := range pokedex.PokemonsCaught {
		out.WriteString(fmt.Sprintf(" - %s\n", key))
	}

	out.WriteString(fmt.Sprintf("Total: %d pokemons\n", len(pokedex.PokemonsCaught)))
	fmt.Print(out.String())
	return nil
}
