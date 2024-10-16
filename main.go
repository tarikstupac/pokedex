package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/tarikstupac/pokedex/internal/pokecache"
	pokedexdata "github.com/tarikstupac/pokedex/internal/pokedex-data"
)

const BLUE = "\033[34m"
const RESET = "\033[0m"
const GREEN = "\033[32m"

var exploreRegexp, _ = regexp.Compile(`^explore\s[\w-]+$`)
var catchRegexp, _ = regexp.Compile(`^catch\s[\w]+$`)
var inspectRegexp, _ = regexp.Compile(`^inspect\s[\w]+$`)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cache := pokecache.NewCache(30 * time.Second)
	pokedex := pokedexdata.NewPokedex()
	confPtr := &config{Next: "https://pokeapi.co/api/v2/location-area/", Previous: "", Cache: cache, Pokedex: pokedex}
	commands := getAvailableCommands()

	for {
		fmt.Print(BLUE + "Pokedex> " + RESET)
		scanner.Scan()
		input := scanner.Text()
		switch input {
		case "help":
			commands[input].callback(&config{}, nil)
		case "exit":
			commands[input].callback(&config{}, nil)
		case "map":
			err := commands[input].callback(confPtr, nil)
			if err != nil {
				fmt.Println(err)
			}
		case "mapb":
			err := commands[input].callback(confPtr, nil)
			if err != nil {
				fmt.Println(err)
			}
		case exploreRegexp.FindString(input):
			if input != "" {
				split := strings.Split(input, " ")
				err := commands[split[0]].callback(confPtr, &split[1])
				if err != nil {
					fmt.Println(err)
				}
			}
		case catchRegexp.FindString(input):
			if input != "" {
				split := strings.Split(input, " ")
				err := commands[split[0]].callback(confPtr, &split[1])
				if err != nil {
					fmt.Println(err)
				}
			}
		case inspectRegexp.FindString(input):
			if input != "" {
				split := strings.Split(input, " ")
				err := commands[split[0]].callback(confPtr, &split[1])
				if err != nil {
					fmt.Println(err)
				}
			}
		default:
			fmt.Printf("%v Command %v unknown, type help for list of commands %v \n", BLUE, input, RESET)
		}
		continue
	}
}
