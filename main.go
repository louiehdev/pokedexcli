package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	var cleanStrings []string
	stringSlice := strings.Fields(text)
	for _, word := range stringSlice {
		cleanStrings = append(cleanStrings, strings.ToLower(word))
	}

	return cleanStrings
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		input := scanner.Text()
		cleanedInput := cleanInput(input)
		if len(cleanedInput) == 0 {
			fmt.Print("ERROR: No input detected!\n")
			continue
		}
		fmt.Printf("Your command was: %v\n", cleanedInput[0])

	}
}
