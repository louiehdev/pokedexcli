package main

import (
	"fmt"
	"bufio"
	"os"
	"time"
	"strings"

	pokecache "github.com/louiehdev/pokedexcli/internal/cache"
	pokeapi "github.com/louiehdev/pokedexcli/internal/utils"
)

type cliCommand struct {
	name        string
	description string
	callback    commandCallback
}

type commandCallback func(client *pokeapi.PokeClient) error

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
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}
	supportedCommands["help"] = cliCommand{
		name:        "help",
		description: "Show list and usage of commands",
		callback:    commandHelp,
	}
	supportedCommands["map"] = cliCommand{
		name:        "map",
		description: "Show list of 20 locations in the Pokemon world",
		callback:    commandMap,
	}
	supportedCommands["mapb"] = cliCommand{
		name:        "mapb",
		description: "Show previous list of 20 locations if available",
		callback:    commandMapb,
	}
	return supportedCommands
}

func commandExit(_client *pokeapi.PokeClient) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(_client *pokeapi.PokeClient) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	commandRegistry := registerCommands()
	for _, command := range commandRegistry {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(client *pokeapi.PokeClient) error {
	url := client.Config.Next
	if url == nil {
		url = &client.BaseURL
	}
	pokeLocationData, err := client.GetPokeLocationData(*url)
	if err != nil {
		return err
	}
	for _, location := range pokeLocationData.Results {
		fmt.Printf("%v\n", location.Name)
	}
	client.Config.Next = pokeLocationData.Next
	client.Config.Previous = pokeLocationData.Previous
	return nil
}

func commandMapb(client *pokeapi.PokeClient) error {
	url := client.Config.Previous
	if url == nil {
		fmt.Print("You are on the first page\n")
		return nil
	}
	pokeLocationData, err := client.GetPokeLocationData(*url)
	if err != nil {
		return err
	}
	for _, location := range pokeLocationData.Results {
		fmt.Printf("%v\n", location.Name)
	}
	client.Config.Next = pokeLocationData.Next
	client.Config.Previous = pokeLocationData.Previous
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
			command.callback(&client)
		} else {
			fmt.Print("Unknown command\n")
			continue
		}
	}
}
