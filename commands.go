package main

import (
	"fmt"
	"os"

	"github.com/tarikstupac/pokedex/internal/pokeapi"
	"github.com/tarikstupac/pokedex/internal/pokecache"
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
			description: "Takes input parameter <area> and displays Pokemon that can be encountered in the area.\n Example usage: explore pastoria-city-area",
			callback:    commandExplore,
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
