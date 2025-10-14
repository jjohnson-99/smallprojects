# Pokedexcli

This is a small cli tool which pulls data from https://pokeapi.co/docs/v2. It can run the following commands, but that's about it! Fun, right?

## Commands

- catch
- exit
- explore
- help
- inspect
- map
- mapb
- pokedex

## The catch command
Calling `catch` with a pokemon's name as the second argument will attempt to catch that pokemon. There is a chance of catching the pokemon depending on the pokemon's BaseExperience as given by the api. If the pokemon is caught, it will be added to your pokedex, which can be viewed by running the `pokedex` command.

## The exit command
Will exit the program, duh.
## The explore command
The `explore` command takes a second argument, the name of the location you want to explore. The command will print a list of pokemon which can be found in that area. This command also uses a cache.

## The help command
Prints a list of commands with their descriptions.

## The inspect command
The `inspect` command takes a pokemon as a second argument. If you have caught that, this command will print information about that pokemon, such as it's height, weight, stats, and types.

## The map and mapb commands
The api paginates the location-areas, each page containing 20 locations. Calling `map` will print the next 20 locations, allowing you to move through the pages. Calling `mapb` will print the previous 20 locations and move your "current" page backwards. There is no "current" page set, the commands just allow you to move through the pages by changing next and previous url's. At each, the page just printed will be cached for a set amount of time, so that calling `map` then quickly calling `mapb` will not ask the api for information and will instead use the temporarily stored information.

## The pokedex command
Prints a list of pokemon you have caught.

### Todo
- [ ] Add documentation.
