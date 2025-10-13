package main

import (
    "fmt"
    "log"
    "io"
    "errors"
    "net/http"
    "encoding/json"
    "github.com/jjohnson-99/smallprojects/pokedexcli/internal"
)

//type Pokemon struct {
//    Pokemon struct {
//        Name string `json:"name"`
//    }`json:"Pokemon"`
//}

//type Encounters struct {
//    Encounter []Pokemon `json:"pokemon_encounters"`
//}
type Encounters struct {
    Encounter []struct {
        Pokemon struct {
            Name string `json:"name"`
        }`json:"Pokemon"`
    }`json:"pokemon_encounters"`
}

func commandExplore(config *Config, cache *internal.Cache, pokedex map[string]Pokemon, parameters []string) error {
    if len(parameters) != 1 {
        err := errors.New(fmt.Sprintf("explore takes 2 parameters, '%v' parameters given", len(parameters)+1))      
        return err
    }
    
    name := parameters[0]
    url := "https://pokeapi.co/api/v2/location-area/"
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

    var encounters Encounters
    if err := json.Unmarshal(body, &encounters); err != nil {
        return err
    }

    fmt.Printf("explore %s...\nFound Pokemon:\n", name)
    for _, encounter := range encounters.Encounter {
        fmt.Println(" - ",encounter.Pokemon.Name)
    }
    
    return nil
}
