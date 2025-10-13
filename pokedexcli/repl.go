package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
    "time"
    "github.com/jjohnson-99/smallprojects/pokedexcli/internal"
)

func startRepl() {
    const interval = 10 * time.Second
    cache := internal.NewCache(interval)
    pokedex := make(map[string]Pokemon)

    config := Config{Next: "https://pokeapi.co/api/v2/location-area/", Previous: ""}
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		reader.Scan()

		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
        parameters := words[1:]

		//command, exists := getCommands(&config, &cache, pokedex, parameters)[commandName]
		command, exists := getCommands()[commandName]
		if exists {
			err := command.callback(&config, &cache, pokedex, parameters)
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, *internal.Cache, map[string]Pokemon, []string) error
}

type Config struct {
    Next string
    Previous string
}

//func getCommands(config *Config, cache *internal.Cache, pokedex map[string]Pokemon, parameters []string) map[string]cliCommand {
	
func getCommands() map[string]cliCommand {
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
			description: "Display the next 20 location area names",
            callback:    commandMap,
		},
	    "mapb": {
			name:        "mapb",
			description: "Display the next 20 location area names",
            callback:    commandMapb,
		},
        "explore" : {
            name:        "explore",
            description: "Display the pokemon in this area",
            callback:    commandExplore,
        },
        "catch": {
            name:        "catch",
            description: "Attempt to catch a pokemon",
            callback:    commandCatch,
        },
        "inspect": {
            name:        "inspect",
            description: "Display details of caught pokemon",
            callback:    commandInspect,
        },
        "pokedex": {
            name:        "pokedex",
            description: "Display all pokemon that have been caught",
            callback:    commandPokedex,
        },
	}
}
