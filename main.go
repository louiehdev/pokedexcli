package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/louiehdev/pokedexcli/internal/utils"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
	config *commandConfig
}

type commandConfig struct {
	Next *string
	Previous *string
}

var Config commandConfig = commandConfig{Next: nil, Previous: nil}

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
		name: "map",
		description: "Show list of 20 locations in the Pokemon world",
		callback: commandMap,
	}
	supportedCommands["mapb"] = cliCommand{
		name: "map",
		description: "Show previous list of 20 locations if available",
		callback: commandMapb,
	}
	return supportedCommands
}

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	commandRegistry := registerCommands()
	for _, command := range commandRegistry {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap() error {
	url := "https://pokeapi.co/api/v2/location-area"
	if Config.Next != nil {
		url = *Config.Next
	}
	pokeLocationData, err := utils.GetPokeLocationData(url)
	if err != nil {
		return err
	}
	for _, location := range pokeLocationData.Results {
		fmt.Printf("%v\n", location.Name)
	}
	Config.Next = pokeLocationData.Next
	Config.Previous = pokeLocationData.Previous
	return nil
}

func commandMapb() error {
	if Config.Previous == nil {
		fmt.Print("You are on the first page\n")
		return nil
	}
	pokeLocationData, err := utils.GetPokeLocationData(*Config.Previous)
	if err != nil {
		return err
	}
	for _, location := range pokeLocationData.Results {
		fmt.Printf("%v\n", location.Name)
	}
	Config.Next = pokeLocationData.Next
	Config.Previous = pokeLocationData.Previous
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commandRegistry := registerCommands()

	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		input := scanner.Text()
		cleanedInput := cleanInput(input)

		command, exists := commandRegistry[cleanedInput[0]]
		if exists {
			command.callback()
		} else {
			fmt.Print("Unknown command\n")
			continue
		}
	}
}
