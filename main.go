package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

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
	return supportedCommands
}

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!")
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
