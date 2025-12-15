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
	callback    func(*Config, []string) error
}

// Создаем переменную для всех команд
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
		"explore": {
			name: "explore",
			description: "See a list of all the Pokémon located in location",
			callback: commandExplore,
		},
	}
	
	// Создаем сканнер для чтения ввода
	scanner := bufio.NewScanner(os.Stdin)
	
	// Основной бесконечный цикл REPL
	for {
		fmt.Print("Pokedex > ")

		if scanner.Scan() {
			// читаем ввод, если пустой перезапуск цикла
			userInput := cleanInput(scanner.Text())
			if len(userInput) == 0 {
				fmt.Println("Type a command, please. 'Help' to see available commands.")
				continue
			}

			// вычленяем команду и проверяем наличие параметров
			commandName := userInput[0]
			arguments := []string{}
			if len(userInput) > 1 {
				arguments = userInput[1:] // Берем все, что после команды
			}

			// Ищем команду и вызываем
			if cmd, exists := commands[commandName]; exists {
				if err := cmd.callback(cfg, arguments); err != nil {
					fmt.Fprintf(os.Stderr, "Error executing %q command: %v\n", commandName, err)
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