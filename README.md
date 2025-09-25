# Pokédex CLI REPL

A REPL that acts as a Pokédex, built in Go!

I created an interactive CLI that uses the PokéAPI to fetch data about Pokémon and locations in the Pokémon world.

It was a fun way to practice Go while learning about:
- HTTP networking
- JSON serialization/deserialization
- Concurrency and caching
- Building interactive CLI tools

## Features

- Explore locations in the Pokémon world
- Encounter and attempt to catch wild Pokémon
- Inspect stats of Pokémon you’ve caught
- Keep track of your caught Pokémon through a Pokédex

## Commands
| Command |	Description |
| :--- | :--- |
|exit|	Exit the Pokédex|
|help|	Show list and usage of commands|
|map|	Show next 20 locations in the Pokémon world|
|mapb|	Show previous 20 locations (if available)|
|explore|	Show Pokémon encountered in a specific location<br />Usage: explore 'location'|
|catch|	Attempt to catch a Pokémon<br />Usage: catch 'pokemon'|
|inspect|	Get information about a caught Pokémon<br />Usage: inspect 'pokemon'|
|pokedex|	Show caught Pokémon in your Pokédex|
