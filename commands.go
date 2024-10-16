package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/tarikstupac/pokedex/internal/pokeapi"
	"github.com/tarikstupac/pokedex/internal/pokecache"
	pokedexdata "github.com/tarikstupac/pokedex/internal/pokedex-data"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, *string) error
}

type config struct {
	Next     string
	Previous string
	Cache    *pokecache.Cache
	Pokedex  *pokedexdata.Pokedex
}

func getAvailableCommands() map[string]cliCommand {
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
		"map": {
			name:        "map",
			description: "Displays 20 locations, subsequent calls display next 20 locations until there's none left",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 locations or errors if there are no previous locations.",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Takes input parameter <area> and displays Pokemon that can be encountered in the area. Example usage: explore pastoria-city-area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Takes input parameter <pokemon> and attempts to catch it. This action can fail. Example usage: catch pikachu.",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Takes input parameter <pokemon> and attempts to inspect it. You can only inspect pokemon you caught. Example usage: inspect pikachu.",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays all the pokemon you've caught so far.",
			callback:    commandPokedex,
		},
	}
}

func commandHelp(conf *config, param *string) error {
	commands := getAvailableCommands()
	fmt.Println(BLUE, "Welcome to the Pokedex!", RESET)
	fmt.Println(BLUE, "Usage:", RESET)
	for _, v := range commands {
		fmt.Printf("%v %v : %v %v\n", GREEN, v.name, v.description, RESET)
	}
	return nil
}

func commandExit(conf *config, param *string) error {
	os.Exit(0)
	return nil
}

func commandMap(conf *config, param *string) error {
	if conf.Next == "" {
		return fmt.Errorf(BLUE + "error: No more locations to display!" + RESET)
	}
	locs, err := pokeapi.GetLocations(conf.Next, conf.Cache)
	if err != nil {
		return fmt.Errorf("error getting locations: %w", err)
	}

	for _, loc := range locs.Results {
		fmt.Println(BLUE, loc.Name, RESET)
	}

	if locs.Next != nil {
		conf.Next = *locs.Next
	} else {
		conf.Next = ""
	}
	if locs.Previous != nil {
		conf.Previous = *locs.Previous
	}

	return nil
}

func commandMapb(conf *config, param *string) error {
	if conf.Previous == "" {
		return fmt.Errorf(BLUE + "Can't go further back" + RESET)
	}
	locs, err := pokeapi.GetLocations(conf.Previous, conf.Cache)
	if err != nil {
		return fmt.Errorf("error getting locations: %w", err)
	}

	for _, loc := range locs.Results {
		fmt.Println(BLUE, loc.Name, RESET)
	}

	if locs.Next != nil {
		conf.Next = *locs.Next
	}
	if locs.Previous != nil {
		conf.Previous = *locs.Previous
	} else {
		conf.Previous = ""
	}
	return nil
}

func commandExplore(conf *config, param *string) error {
	if param == nil {
		return fmt.Errorf("error empty value supplied for param")
	}
	locDetails, err := pokeapi.GetEncounters(*param, conf.Cache)
	if err != nil {
		return fmt.Errorf("error getting encounters: %w", err)
	}
	fmt.Println(BLUE, "Exploring ", *param, " area...", RESET)
	fmt.Println(BLUE, "Found pokemon: ", RESET)
	for _, val := range locDetails.PokemonEncounters {
		fmt.Println(BLUE, "- ", val.Pokemon.Name, RESET)
	}
	return nil
}

func commandCatch(conf *config, param *string) error {
	if param == nil {
		return fmt.Errorf("error empty value supplied for param")
	}
	pokemon, err := pokeapi.GetPokemon(*param, conf.Cache)
	if err != nil {
		return fmt.Errorf("error getting pokemon: %w", err)
	}

	fmt.Println(BLUE, "Throwing a pokeball at ", *param, " ...", RESET)
	catchVal := rand.Intn(pokemon.BaseExperience * 2)
	if catchVal >= pokemon.BaseExperience {
		conf.Pokedex.Add(*param, pokemon)

		fmt.Println(BLUE, *param, " was caught!", RESET)
	} else {
		fmt.Println(BLUE, *param, " escaped!", RESET)
	}

	return nil
}

func commandInspect(conf *config, param *string) error {
	if param == nil {
		return fmt.Errorf("error empty value supplied for param")
	}
	pokemon, ok := conf.Pokedex.Get(*param)
	if !ok {
		fmt.Println(BLUE, "you have not caught: ", *param, " or pokemon with that name doesn't exit!", RESET)
		return nil
	}

	fmt.Println(BLUE, "Name: ", pokemon.Name, RESET)
	fmt.Println(BLUE, "Height: ", pokemon.Height, RESET)
	fmt.Println(BLUE, "Weight: ", pokemon.Weight, RESET)
	fmt.Println(BLUE, "Stats: ", RESET)
	for _, stat := range pokemon.Stats {
		fmt.Println(BLUE, " -", stat.Stat.Name, ": ", stat.BaseStat, RESET)
	}
	fmt.Println(BLUE, "Types: ", RESET)
	for _, t := range pokemon.Types {
		fmt.Println(BLUE, " -", t.Type.Name, RESET)
	}

	return nil

}

func commandPokedex(conf *config, param *string) error {
	caughtPokemon := conf.Pokedex.GetCaughtPokemon()
	if len(caughtPokemon) == 0 {
		fmt.Println(BLUE, "You haven't caught any pokemon yet, type help to explore and catch some pokemon!", RESET)
		return nil
	}
	fmt.Println(BLUE, "Your Pokedex: ", RESET)
	for key := range caughtPokemon {
		fmt.Println(BLUE, " -", key, RESET)
	}
	return nil
}
