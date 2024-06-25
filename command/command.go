package command

import (
	"fmt"
	"os"

	"github.com/pakkermandev/go-pokedex/api"
)

type CliCommand struct {
	Name        string
	Description string
	Callback    func() error
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
	}
}

func commandHelp() error {
	fmt.Println("\nUsage:")

	options := GetOptions()

	for _, opt := range options {
		fmt.Printf("%v:\t%v\n", opt.Name, opt.Description)
	}

	fmt.Println()
	return nil
}

func commandExit() error {
	fmt.Println("Exit the program")
	os.Exit(0)
	return nil
}

func mapLocations() error {
	api.GetNextMap()

	return nil
}

func mapbLocations() error {
	api.GetPreviousMap()
	return nil
}
