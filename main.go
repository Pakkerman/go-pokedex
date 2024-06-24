package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	options := getOptions()

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		fmt.Println("you entered: ", input)

		option := options[input]
		option.callback()
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func getOptions() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}

func commandHelp() error {
	fmt.Println("Welcome to Pokedex! \nUsage: \n")

	options := getOptions()

	for _, opt := range options {
		fmt.Printf("%v: %v\n", opt.name, opt.description)
	}

	return nil
}

func commandExit() error {
	fmt.Println("Exit the program")
	os.Exit(0)
	return nil
}
