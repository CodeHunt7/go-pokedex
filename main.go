package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("Pokedex > ")

		if scanner.Scan() {
			userInput := scanner.Text()
			firstWord := cleanInput(userInput)[0]
			fmt.Printf("Your command was: %s\n", firstWord)
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
		}	
	}
	
}
