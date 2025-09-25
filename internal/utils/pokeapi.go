package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"

	pokecache "github.com/louiehdev/pokedexcli/internal/cache"
)

type CommonData struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokeAreaData struct {
	Next      *string
	Previous  *string
	Locations []CommonData `json:"results"`
}

type PokeLocationData struct {
	Location          CommonData
	PokemonEncounters []PokeEncounter `json:"pokemon_encounters"`
}

type PokeEncounter struct {
	Pokemon CommonData
}

type PokemonData struct {
	Id             int        `json:"id"`
	Name           string     `json:"name"`
	BaseExperience int        `json:"base_experience"`
	Height         int        `json:"height"`
	Weight         int        `json:"weight"`
	Stats          []PokeStat `json:"stats"`
	Types          []PokeType `json:"types"`
}

type PokeStat struct {
	Stat  CommonData `json:"stat"`
	Value int        `json:"base_stat"`
}

type PokeType struct {
	Type CommonData `json:"type"`
}

type PokeClient struct {
	httpClient *http.Client
	cache      *pokecache.Cache
	BaseURL    string
	Config     APIConfig
	Pokedex    Pokedex
}

type APIConfig struct {
	Next     *string
	Previous *string
}

func (p *PokeClient) GetPokeAreaData(url string) (PokeAreaData, error) {
	var areaData PokeAreaData

	if entry, exists := p.cache.Get(url); exists {
		json.Unmarshal(entry, &areaData)
		return areaData, nil
	}

	res, err := p.httpClient.Get(url)
	if err != nil {
		return PokeAreaData{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return PokeAreaData{}, err
	}
	p.cache.Add(url, data)

	if err := json.Unmarshal(data, &areaData); err != nil {
		return PokeAreaData{}, err
	}

	return areaData, nil
}

func (p *PokeClient) GetPokeLocationData(location string) (PokeLocationData, error) {
	var locationData PokeLocationData

	url := p.BaseURL + "location-area/" + location
	if entry, exists := p.cache.Get(url); exists {
		json.Unmarshal(entry, &locationData)
		return locationData, nil
	}

	res, err := p.httpClient.Get(url)
	if err != nil {
		return PokeLocationData{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return PokeLocationData{}, err
	}
	p.cache.Add(url, data)

	if err := json.Unmarshal(data, &locationData); err != nil {
		return PokeLocationData{}, err
	}

	return locationData, nil
}

func (p *PokeClient) GetPokemonData(pokename string) (PokemonData, error) {
	var pokemonData PokemonData

	url := p.BaseURL + "pokemon/" + pokename
	if entry, exists := p.cache.Get(url); exists {
		json.Unmarshal(entry, &pokemonData)
		return pokemonData, nil
	}
	res, err := p.httpClient.Get(url)
	if err != nil {
		return PokemonData{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return PokemonData{}, err
	}
	p.cache.Add(url, data)

	if err := json.Unmarshal(data, &pokemonData); err != nil {
		return PokemonData{}, err
	}

	return pokemonData, nil
}

type Pokedex struct {
	Data map[string]PokemonData
}

func (d *Pokedex) Add(key string, data PokemonData) {
	d.Data[key] = data
}

func (d *Pokedex) Get(key string) (PokemonData, bool) {
	pokemon, exists := d.Data[key]
	if exists {
		return pokemon, true
	} else {
		return PokemonData{}, false
	}
}

func NewClient(cache *pokecache.Cache) PokeClient {
	var client PokeClient
	client.httpClient = &http.Client{}
	client.cache = cache
	client.BaseURL = "https://pokeapi.co/api/v2/"
	client.Config = APIConfig{}
	client.Pokedex = Pokedex{Data: make(map[string]PokemonData)}
	return client
}
