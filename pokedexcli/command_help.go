package main

import (
    "fmt"
    "github.com/jjohnson-99/smallprojects/pokedexcli/internal"
)

func commandHelp(config *Config, cache *internal.Cache, pokedex map[string]Pokemon, parameters []string) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	//for _, cmd := range getCommands(config, cache, pokedex, parameters) {
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()
	return nil
}
