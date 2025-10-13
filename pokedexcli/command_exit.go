package main

import (
	"fmt"
	"os"
    "github.com/jjohnson-99/smallprojects/pokedexcli/internal"
)

func commandExit(config *Config, cache *internal.Cache, pokemon map[string]Pokemon, parameters []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
