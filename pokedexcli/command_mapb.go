package main

import (
	"fmt"
    "io"
    "log"
    "net/http"
    "encoding/json"
    "errors"
    "github.com/jjohnson-99/smallprojects/pokedexcli/internal"
)

func commandMapb(config *Config, cache *internal.Cache, pokedex map[string]Pokemon, parameters []string) error {
    if config.Previous == "" {
        err := errors.New("There is no previous page of location areas")
        return err
    }
    url := config.Previous
    body, ok := cache.Get(url)
    if !ok {
        res, err := http.Get(url)
        if err != nil {
            log.Fatal(err)
        }
        body, err = io.ReadAll(res.Body)
        res.Body.Close()
        if res.StatusCode > 299 {
            log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
        }
        if err != nil {
            log.Fatal(err)
        }
    }

    var page internal.LocationAreaPage
    if err := json.Unmarshal(body, &page); err != nil {
        return err
    }

    results := page.Results
    for _, location := range results {
        fmt.Println(location.Name)
    }

    config.Next = *page.Next
    if page.Previous != nil {
        config.Previous = *page.Previous
    }

	return nil
}
