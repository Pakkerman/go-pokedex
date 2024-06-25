package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pakkermandev/go-pokedex/command"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	options := command.GetOptions()

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()

		// utils.ClearConsole()

		option, ok := options[input]
		if !ok {
			fmt.Println("error: unknown command")
			options["help"].Callback()
			continue
		}

		option.Callback()

	}
}
