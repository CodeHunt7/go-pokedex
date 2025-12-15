package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/CodeHunt7/go-pokedex/internal/pokecache"
)

// Структура для хранения состояния ссылок на API
type Config struct {
    Next      string
    Previous  string
    pokeCache pokecache.Cache
}

// Структура для распаковки JSON ответа от PokeAPI по списку локаций
type LocationAreaResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// Структура для распаковки JSON ответа от PokeAPI по конкретной локации
type ConcreteLocationResponce struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int           `json:"chance"`
				ConditionValues []interface{} `json:"condition_values"`
				MaxLevel        int           `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
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

func commandExit(cfg *Config, parameters []string) error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}

func commandHelp(cfg *Config, parameters []string) error {
    fmt.Println("\nWelcome to the Pokedex!")
    fmt.Println("Usage:")
    for _, cmd := range commands {
        fmt.Printf("%s: %s\n", cmd.name, cmd.description)
    }
    fmt.Println()
    return nil
}

func commandMap(cfg *Config, parameters []string) error {
    // Инициализируем ссылки
    NextURL := cfg.Next
    if NextURL == "" { 
        NextURL = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
    }

    // Создаем переменные для тела ответа и ошибки
    var body []byte
    var err error

    // Проверяем есть ли ответ в кеше
    cachedResponse, inCache := cfg.pokeCache.Get(NextURL)
    
    if !inCache { // В кеше нет, делаем запрос
        // Получаем JSON ответ от PokeAPI
        res, err := http.Get(NextURL)
        if err != nil { 
            return err
        }
        body, err = io.ReadAll(res.Body)
        res.Body.Close()
        if res.StatusCode > 299 {
            return fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
        }
        if err != nil {
            return err
        }
        //fmt.Println(body)
        //fmt.Println("FROM API!")
    } else { // В кеше есть, используем его
        body = cachedResponse
        //fmt.Println("FROM CACHE!")
    }
    
    // Распаковываем JSON в стуктуру
    var locations LocationAreaResponse
    err = json.Unmarshal(body, &locations)
    if err != nil {
        return err
    }

    // Обновляем глобальный конфиг и кеш
    cfg.Next = locations.Next
    if locations.Previous != nil {
        cfg.Previous = locations.Previous.(string)
    } else {
        cfg.Previous = ""
    }
    cfg.pokeCache.Add(NextURL, body)

    // Выводим ответ в консоль
    fmt.Println()
    for _, location := range locations.Results {
        fmt.Printf(" - %s\n", location.Name)
    }
    fmt.Println()

    return nil
}

func commandMapBack(cfg *Config, parameters []string) error {
    // Инициализируем ссылки
    PrevURL := cfg.Previous
    if PrevURL == "" {
        fmt.Println("No previous locations available.")
        return nil
    }

    // Создаем переменные для тела ответа и ошибки
    var body []byte
    var err error

    // Проверяем есть ли в кеше
    cachedResponse, inCache := cfg.pokeCache.Get(PrevURL)
    
    if !inCache { // В кеше нет, делаем запрос
        // Получаем JSON ответ от PokeAPI
        res, err := http.Get(PrevURL)
        if err != nil { 
            return err
        }
        body, err = io.ReadAll(res.Body)
        res.Body.Close()
        if res.StatusCode > 299 {
            return fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
        }
        if err != nil {
            return err
        }
        // fmt.Println("FROM API!")
        // fmt.Println(PrevURL)
    } else { // В кеще есть, используем его
        body = cachedResponse
        // fmt.Println("FROM CACHE!")
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
    cfg.pokeCache.Add(PrevURL, body)

    // Выводим ответ в консоль
    fmt.Println()
    for _, location := range locations.Results {
        fmt.Printf(" - %s\n", location.Name)
    }
    fmt.Println()

    return nil
}

func commandExplore(cfg *Config, parameters []string) error {
    // инициализируем ссылку
    if len(parameters) == 0 {
        fmt.Println("Please provide a location area name to explore.")
        return nil
    }
    exploreURL := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", parameters[0])

    // Создаем переменные для тела ответа и ошибки
    var body []byte
    var err error

    // проверяем есть ли в кеше
    cachedResponse, inCache := cfg.pokeCache.Get(exploreURL)
    if !inCache { // В кеше нет, делаем запрос
        // Получаем JSON ответ от PokeAPI
        res, err := http.Get(exploreURL)
        if err != nil { 
            return err
        }
        body, err = io.ReadAll(res.Body)
        res.Body.Close()
        if res.StatusCode > 299 {
            return fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
        }
        if err != nil {
            return err
        }
    } else { // В кеще есть, используем его
        body = cachedResponse
    }

    // Распаковываем JSON в стуктуру
    var locationInfo ConcreteLocationResponce
    err = json.Unmarshal(body, &locationInfo)
    if err != nil {
        return err
    }

    // Обновляем глобальный кеш в конфиге
    cfg.pokeCache.Add(exploreURL, body)

    // Выводим ответ в консоль
    fmt.Println()
    fmt.Printf("Exploring %s...\n", parameters[0])
    fmt.Println("Found Pokemon:")
    for _, pokemonInfo := range locationInfo.PokemonEncounters {
        fmt.Printf(" - %s\n", pokemonInfo.Pokemon.Name)
    }
    fmt.Println()

    return nil
}