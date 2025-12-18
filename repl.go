package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/CodeHunt7/go-pokedex/internal/pokecache"
)

// Константа для сложности поимки покемона
const ThrowingDifficulty = 40

// Структура для хранения состояния ссылок на API
type Config struct {
    Next      string
    Previous  string
    pokeCache pokecache.Cache
    Pokedex   map[string]PokemonResponse
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

// Структура для распаковки JSON ответа от PokeAPI по конкретному покемону
type PokemonResponse struct {
	Abilities []struct {
		Ability struct {   
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height    int `json:"height"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			Order        interface{} `json:"order"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []struct {
		Abilities []struct {
			Ability  interface{} `json:"ability"`
			IsHidden bool        `json:"is_hidden"`
			Slot     int         `json:"slot"`
		} `json:"abilities"`
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
	} `json:"past_abilities"`
	PastTypes []interface{} `json:"past_types"`
	Species   struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string `json:"back_default"`
		BackFemale       string `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  string `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      string `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale string `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string      `json:"front_default"`
				FrontFemale  interface{} `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      string `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale string `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string      `json:"back_default"`
				BackFemale       string      `json:"back_female"`
				BackShiny        string      `json:"back_shiny"`
				BackShinyFemale  interface{} `json:"back_shiny_female"`
				FrontDefault     string      `json:"front_default"`
				FrontFemale      string      `json:"front_female"`
				FrontShiny       string      `json:"front_shiny"`
				FrontShinyFemale string      `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault           string `json:"back_default"`
					BackShiny             string `json:"back_shiny"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					BackTransparent       string `json:"back_transparent"`
					FrontDefault          string `json:"front_default"`
					FrontShiny            string `json:"front_shiny"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
					FrontTransparent      string `json:"front_transparent"`
				} `json:"crystal"`
				Gold struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"gold"`
				Silver struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationIx struct {
				ScarletViolet struct {
					FrontDefault string      `json:"front_default"`
					FrontFemale  interface{} `json:"front_female"`
				} `json:"scarlet-violet"`
			} `json:"generation-ix"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default"`
						BackFemale       string `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  string `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      string `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale string `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string      `json:"front_default"`
					FrontFemale  interface{} `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				BrilliantDiamondShiningPearl struct {
					FrontDefault string      `json:"front_default"`
					FrontFemale  interface{} `json:"front_female"`
				} `json:"brilliant-diamond-shining-pearl"`
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  string `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
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

func commandCatch(cfg *Config, parameters []string) error {
    // инициализируем ссылку
    if len(parameters) == 0 {
        fmt.Println("Please provide a name of the pokemon to catch.")
        return nil
    }
    catchURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", parameters[0])

    // Создаем переменные для тела ответа и ошибки
    var body []byte
    var err error

    // проверяем есть ли в кеше
    cachedResponse, inCache := cfg.pokeCache.Get(catchURL)
    if !inCache { // В кеше нет, делаем запрос
        // Получаем JSON ответ от PokeAPI
        res, err := http.Get(catchURL)
        if err != nil { 
            return err
        }
        body, err = io.ReadAll(res.Body)
        res.Body.Close()

        // поверяем, что покемон существует
        if res.StatusCode == 404 {
            fmt.Printf("%s is not a valid pokemon name\n", parameters[0])
            return nil
        }

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
    var pokemonInfo PokemonResponse
    err = json.Unmarshal(body, &pokemonInfo)
    if err != nil {
        return err
    }

    // Считаем шансы поймать и Выводим ответ в консоль
    fmt.Println()
    fmt.Printf("Throwing a Pokeball at %s...\n", parameters[0])
    throwResult := rand.Intn(pokemonInfo.BaseExperience)
    
    if throwResult < ThrowingDifficulty { // поймал
        fmt.Printf("%s was caught!\n", pokemonInfo.Name)
        cfg.Pokedex[pokemonInfo.Name] = pokemonInfo
        //fmt.Printf("Deb: \n%d - base exp \n%d - throw res \n%d - defficulty\n", pokemonInfo.BaseExperience, throwResult, ThrowingDifficulty)
    } else { // не поймал
        fmt.Printf("%s escaped!\n", pokemonInfo.Name)
        //fmt.Printf("Deb: \n%d - base exp \n%d - throw res \n%d - defficulty\n", pokemonInfo.BaseExperience, throwResult, ThrowingDifficulty)
    }

    // Обновляем глобальный кеш в конфиге
    cfg.pokeCache.Add(catchURL, body)

    return nil
}