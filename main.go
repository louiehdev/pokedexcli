package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	pokecache "github.com/louiehdev/pokedexcli/internal/cache"
	pokeapi "github.com/louiehdev/pokedexcli/internal/utils"
)

type cliCommand struct {
	name         string
	description  string
	callback     commandCallback
	requireInput bool
}

type commandCallback func(client *pokeapi.PokeClient, input string) error

func cleanInput(text string) []string {
	var cleanStrings []string
	stringSlice := strings.Fields(text)
	for _, word := range stringSlice {
		cleanStrings = append(cleanStrings, strings.ToLower(word))
	}

	return cleanStrings
}

func registerCommands() map[string]cliCommand {
	supportedCommands := make(map[string]cliCommand)
	supportedCommands["exit"] = cliCommand{
		name:         "exit",
		description:  "Exit the Pokedex",
		callback:     commandExit,
		requireInput: false,
	}
	supportedCommands["help"] = cliCommand{
		name:         "help",
		description:  "Show list and usage of commands",
		callback:     commandHelp,
		requireInput: false,
	}
	supportedCommands["map"] = cliCommand{
		name:         "map",
		description:  "Show list of 20 locations in the Pokemon world",
		callback:     commandMap,
		requireInput: false,
	}
	supportedCommands["mapb"] = cliCommand{
		name:         "mapb",
		description:  "Show previous list of 20 locations if available",
		callback:     commandMapb,
		requireInput: false,
	}
	supportedCommands["explore"] = cliCommand{
		name:         "explore",
		description:  "Show Pokemon encountered in a specific location",
		callback:     commandExplore,
		requireInput: true,
	}
	supportedCommands["catch"] = cliCommand{
		name:         "catch",
		description:  "Attempt to catch a Pokemon",
		callback:     commandCatch,
		requireInput: true,
	}
	supportedCommands["inspect"] = cliCommand{
		name:         "inspect",
		description:  "Get information about a caught Pokemon",
		callback:     commandInspect,
		requireInput: true,
	}
	supportedCommands["pokedex"] = cliCommand{
		name:         "pokedex",
		description:  "Shows caught Pokemon in your Pokedex",
		callback:     commandPokedex,
		requireInput: false,
	}
	return supportedCommands
}

func commandExit(_client *pokeapi.PokeClient, _input string) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(_client *pokeapi.PokeClient, _input string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	commandRegistry := registerCommands()
	for _, command := range commandRegistry {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(client *pokeapi.PokeClient, _input string) error {
	url := client.Config.Next
	if url == nil {
		baseurl := client.BaseURL + "location-area/"
		url = &baseurl
	}
	pokeAreaData, err := client.GetPokeAreaData(*url)
	if err != nil {
		return err
	}
	for _, location := range pokeAreaData.Locations {
		fmt.Printf("%v\n", location.Name)
	}
	client.Config.Next = pokeAreaData.Next
	client.Config.Previous = pokeAreaData.Previous
	return nil
}

func commandMapb(client *pokeapi.PokeClient, _input string) error {
	url := client.Config.Previous
	if url == nil {
		fmt.Print("You are on the first page\n")
		return nil
	}
	pokeAreaData, err := client.GetPokeAreaData(*url)
	if err != nil {
		return err
	}
	for _, location := range pokeAreaData.Locations {
		fmt.Printf("%v\n", location.Name)
	}
	client.Config.Next = pokeAreaData.Next
	client.Config.Previous = pokeAreaData.Previous
	return nil
}

func commandExplore(client *pokeapi.PokeClient, location string) error {
	pokeLocationData, err := client.GetPokeLocationData(location)
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %s...\n", location)
	fmt.Print("Pokemon encountered:\n")
	for _, encounter := range pokeLocationData.PokemonEncounters {
		fmt.Printf(" - %v\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(client *pokeapi.PokeClient, pokemon string) error {
	pokemonData, err := client.GetPokemonData(pokemon)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonData.Name)
	successChance := min(1.0, 75.0/float64(pokemonData.BaseExperience))
	randomChance := rand.Float64()

	if randomChance < successChance {
		fmt.Printf("%s was caught! It can now be inspected.\n", pokemonData.Name)
		client.Pokedex.Add(pokemon, pokemonData)
		return nil
	} else {
		fmt.Printf("%s escaped!\n", pokemonData.Name)
		return nil
	}
}

func commandInspect(client *pokeapi.PokeClient, pokemon string) error {
	if p, exists := client.Pokedex.Get(pokemon); exists {
		fmt.Printf("Name: %s\nHeight: %v\nWeight: %v\n", p.Name, p.Height, p.Weight)
		fmt.Println("Stats:")
		for _, s := range p.Stats {
			fmt.Printf(" - %v: %v\n", s.Stat.Name, s.Value)
		}
		fmt.Println("Types:")
		for _, t := range p.Types {
			fmt.Printf(" - %v\n", t.Type.Name)
		}
		return nil
	} else {
		fmt.Println("You have not caught that pokemon yet")
		return nil
	}
}

func commandPokedex(client *pokeapi.PokeClient, _input string) error {
	pokedexData := client.Pokedex.Data
	fmt.Println("Your Pokedex:")
	for pokemon := range pokedexData {
		fmt.Printf(" - %s\n", pokemon)
	}
	return nil
}

func main() {
	newCache := pokecache.NewCache(5 * time.Second)
	client := pokeapi.NewClient(newCache)
	scanner := bufio.NewScanner(os.Stdin)
	commandRegistry := registerCommands()

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleanedInput := cleanInput(input)

		command, exists := commandRegistry[cleanedInput[0]]
		if exists {
			if len(cleanedInput) < 2 && command.requireInput {
				fmt.Print("Command missing args\n")
				continue
			} else if len(cleanedInput) < 2 {
				command.callback(&client, "")
			} else {
				command.callback(&client, cleanedInput[1])
			}
		} else {
			fmt.Print("Unknown command\n")
			continue
		}
	}
}
