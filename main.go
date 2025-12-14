package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/CodeHunt7/go-pokedex/internal/pokecache"
)

// Структура для CLI команд
type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

var commands map[string]cliCommand

func main() {
	// Делаем конфиг
	cache := pokecache.NewCache(1 * time.Minute)
	cfg := &Config{
		pokeCache: cache,
	}

	// Инициализируем команды
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
		"map": {
			name: "map",
			description: "Displays the names of 20 location areas",
			callback: commandMap,
		},
		"mapb": {
			name: "mapb",
			description: "Displays the names of the previous 20 location areas",
			callback: commandMapBack,
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
				if err := cmd.callback(cfg); err != nil {
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