package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands map[string]cliCommand

func main() {
	
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
	}
	
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("Pokedex > ")

		if scanner.Scan() {
			userInput := scanner.Text()
			firstWord := cleanInput(userInput)[0]
			//fmt.Printf("Your command was: %s\n", firstWord)

			if cmd, exists := commands[firstWord]; exists {
				if err := cmd.callback(); err != nil {
					fmt.Fprintf(os.Stderr, "Error executing command %q: %v\n", firstWord, err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
		}	
	}
	
}