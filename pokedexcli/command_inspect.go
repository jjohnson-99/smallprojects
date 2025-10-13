package main

import (
    "fmt"
    "errors"
    "github.com/jjohnson-99/smallprojects/pokedexcli/internal"
)

/*import (
    "fmt"
    "math/rand"
    "net/http"
    "log"
    "io"
    "errors"
    "encoding/json"
    "github.com/jjohnson-99/pokedexcli/internal"
)

func commandInspect(config *Config, cache *internal.Cache, parameters []string) error {
    if len(parameters) != 1 {
        err := errors.New(fmt.Sprintf("explore takes 2 parameters, '%v' parameters given", len(parameters)+1))      
        return err
    }
    
    name := parameters[0]
    url := "https://pokeapi.co/api/v2/pokemon/"
    url = url + name
    body, ok := cache.Get(url)

    return nil
}*/


func commandInspect(config *Config, cache *internal.Cache, pokedex map[string]Pokemon, parameters []string) error {
    if len(parameters) != 1 {
        err := errors.New(fmt.Sprintf("explore takes 2 parameters, '%v' parameters given", len(parameters)+1))      
        return err
    }   
    
    name := parameters[0]
    pokemon, ok := pokedex[name]
    if !ok {
        fmt.Println("you have not caught that pokemon")
        return nil
    }
    fmt.Printf("Name: %s\n", pokemon.Name)
    fmt.Printf("Height: %d\n", pokemon.Height) 
    fmt.Printf("Weight: %d\n", pokemon.Weight) 
    fmt.Print("Stats:\n")
    for _, stat := range pokemon.Stats {
        fmt.Printf("  -%s: %d\n", stat.Stat.StatName, stat.BaseStat)
    }
    fmt.Print("Types:\n")
    for _, t := range pokemon.Types {
        fmt.Printf("  - %s\n", t.Type.TypeName)
    }
    return nil

}
