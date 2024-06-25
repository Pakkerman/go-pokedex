package command

import (
	"fmt"
	"os"

	"github.com/pakkermandev/go-pokedex/api"
)

type CliCommand struct {
	Name        string
	Description string
	Callback    func(*string) error
}

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
	api.GetNextMap()

	return nil
}

func mapbLocations(arg *string) error {
	api.GetPreviousMap()
	return nil
}

func explore(arg *string) error {
	name := *arg
	api.Explore(name)
	return nil
}
