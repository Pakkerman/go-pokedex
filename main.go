package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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

		input := strings.Split(scanner.Text(), " ")

		if len(input) == 0 {
			continue
		}

		command := input[0]
		var arg *string

		if len(input) > 1 {
			arg = &input[1]
		}

		if cmd, ok := options[command]; ok {
			err := cmd.Callback(arg)
			if err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Command not found:", command)
		}

		fmt.Println("Enter command:")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}
