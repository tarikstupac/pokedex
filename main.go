package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/tarikstupac/pokedex/internal/pokecache"
)

const BLUE = "\033[34m"
const RESET = "\033[0m"
const GREEN = "\033[32m"

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	confPtr := &config{Next: "https://pokeapi.co/api/v2/location/", Previous: ""}
	commands := getAvailableCommands()
	cache := pokecache.NewCache(30 * time.Second)

	for {
		fmt.Print(BLUE + "Pokedex> " + RESET)
		scanner.Scan()
		input := scanner.Text()
		switch input {
		case "help":
			commands[input].callback(&config{}, cache)
		case "exit":
			commands[input].callback(&config{}, cache)
		case "map":
			err := commands[input].callback(confPtr, cache)
			if err != nil {
				fmt.Println(err)
			}
		case "mapb":
			err := commands[input].callback(confPtr, cache)
			if err != nil {
				fmt.Println(err)
			}
		default:
			fmt.Printf("%v Command %v unknown, type help for list of commands %v \n", BLUE, input, RESET)
		}
		continue
	}
}
