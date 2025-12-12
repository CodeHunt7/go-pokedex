package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Структура для хранения состояния ссылок на API
type Config struct {
    Next     string
    Previous string
}

// Структура для распаковки JSON ответа от PokeAPI
type LocationAreaResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func cleanInput(text string) []string {
    // Fields сама разобьет по пробелам и уберет дублирующиеся пробелы
    fields := strings.Fields(text)
    
    // Если слов много, лучше сразу выделить память (оптимизация)
    output := make([]string, 0, len(fields))
    
    for _, word := range fields {
        // Просто лоуеркейсим и добавляем в ответ
        output = append(output, strings.ToLower(word))
    }
    
    return output
}

func commandExit(cfg *Config) error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}

func commandHelp(cfg *Config) error {
    fmt.Println("Welcome to the Pokedex!")
    fmt.Println("Usage:")
    for _, cmd := range commands {
        fmt.Printf("%s: %s\n", cmd.name, cmd.description)
    }
    fmt.Println()
    return nil
}

func commandMap(cfg *Config) error {
    // Инициализируем ссылки
    NextURL := cfg.Next
    if NextURL == "" { 
        NextURL = "https://pokeapi.co/api/v2/location-area/"
    }

    // Получаем JSON ответ от PokeAPI
    res, err := http.Get(NextURL)
    if err != nil { 
        return err
    }
    body, err := io.ReadAll(res.Body)
    res.Body.Close()
    if res.StatusCode > 299 {
		return fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}
    if err != nil {
        return err
    }
    //fmt.Println(body)
    
    // Распаковываем JSON в стуктуру
    var locations LocationAreaResponse
    err = json.Unmarshal(body, &locations)
    if err != nil {
        return err
    }

    // Обновляем глобальный конфиг
    cfg.Next = locations.Next
    if locations.Previous != nil {
        cfg.Previous = locations.Previous.(string)
    } else {
        cfg.Previous = ""
    }

    // Выводим ответ в консоль
    fmt.Println()
    for _, location := range locations.Results {
        fmt.Println(location.Name)
    }
    fmt.Println()

    return nil
}

func commandMapBack(cfg *Config) error {
    
    // Инициализируем ссылки
    PrevURL := cfg.Previous
    if PrevURL == "" {
        fmt.Println("No previous locations available.")
        return nil
    }

    // Получаем JSON ответ от PokeAPI
    res, err := http.Get(PrevURL)
    if err != nil { 
        return err
    }
    body, err := io.ReadAll(res.Body)
    res.Body.Close()
    if res.StatusCode > 299 {
		return fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}
    if err != nil {
        return err
    }

    // Распаковываем JSON в стуктуру
    var locations LocationAreaResponse
    err = json.Unmarshal(body, &locations)
    if err != nil {
        return err
    }

    // Обновляем глобальный конфиг
    cfg.Next = locations.Next
    if locations.Previous != nil {
        cfg.Previous = locations.Previous.(string)
    } else {
        cfg.Previous = ""
    }

    // Выводим ответ в консоль
    fmt.Println()
    for _, location := range locations.Results {
        fmt.Println(location.Name)
    }
    fmt.Println()

    return nil
}