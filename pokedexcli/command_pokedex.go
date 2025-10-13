package main

import (
    "fmt"
    "github.com/jjohnson-99/smallprojects/pokedexcli/internal"
)

func commandPokedex(config *Config, cache *internal.Cache, pokedex map[string]Pokemon, parameters []string) error {
    for name, _ := range pokedex {
        fmt.Printf(" - %s\n", name)
    }
    return nil
}
