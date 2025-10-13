package main

import (
    "fmt"
    "math/rand"
    "net/http"
    "log"
    "io"
    "errors"
    "encoding/json"
    "github.com/jjohnson-99/smallprojects/pokedexcli/internal"
)


type Pokemon struct {
    Name string `json:"name"`
    BaseExperience int `json:"base_experience"`
    Height int `json:"height"`
    Weight int `json:"weight"`
    Stats []struct {
        BaseStat int `json:"base_stat"`
        Stat struct {
            StatName string `json:"name"`
        }`json:"stat"`
    }`json:"stats"`
    Types []struct {
        Type struct {
            TypeName string `json:"name"`
        }`json:"type"`
    }`json:"types"`
}

// Name:
// Height:
// Weight:
// Stats:
//   -hp:
//   -attack:
//   -defense:
//   -special-attack:
//   -special-defense:
//   -speed:
// Types:
//   - normal
//   - flying

func commandCatch(config *Config, cache *internal.Cache, pokedex map[string]Pokemon, parameters []string) error {
    if len(parameters) != 1 {
        err := errors.New(fmt.Sprintf("explore takes 2 parameters, '%v' parameters given", len(parameters)+1))      
        return err
    }
    
    name := parameters[0]
    url := "https://pokeapi.co/api/v2/pokemon/"
    url = url + name
    body, ok := cache.Get(url)
    if !ok {
        res, err := http.Get(url)
        if err != nil {
            log.Fatal(err)
        }
        defer res.Body.Close()
        body, err = io.ReadAll(res.Body)
        cache.Add(url, body)

        if res.StatusCode > 299 {
            log.Fatalf("Response failed with status code: %d and \nbody: %s\n", res.StatusCode, body)
        }
        if err != nil {
            log.Fatal(err)
        }
    }

    var pokemon Pokemon
    if err := json.Unmarshal(body, &pokemon); err != nil {
        return err
    }
    random := rand.Intn(1000)
    fmt.Printf("Throwing a Pokeball at %s...\n", name)
    if random >= pokemon.BaseExperience {
        pokedex[name] = pokemon
        fmt.Printf("%s was caught!\n", name)
    } else {
        fmt.Printf("%s escaped!\n", name)
    }
    return nil
}
