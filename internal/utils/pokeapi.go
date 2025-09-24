package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"

	pokecache "github.com/louiehdev/pokedexcli/internal/cache"
)

type PokeClient struct {
	httpClient *http.Client
	cache      *pokecache.Cache
	BaseURL    string
	Config     APIConfig
	Pokedex map[string]Pokemon
}

type APIConfig struct {
	Next     *string
	Previous *string
}

type PokeData struct {
	Name string
	Url  string
}

type PokeAreaData struct {
	Next      *string
	Previous  *string
	Locations []PokeData `json:"results"`
}

type PokeLocationData struct {
	Location          PokeData
	PokemonEncounters []PokeEncounter `json:"pokemon_encounters"`
}

type PokeEncounter struct {
	Pokemon PokeData
}

type Pokemon struct {
	Id int
	Name string
	BaseExperience int `json:"base_experience"`
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

func (p *PokeClient) GetPokemonData(pokename string) (Pokemon, error) {
	var pokemonData Pokemon

	url := p.BaseURL + "pokemon/" + pokename
	if entry, exists := p.cache.Get(url); exists {
		json.Unmarshal(entry, &pokemonData)
		return pokemonData, nil
	}
	res, err := p.httpClient.Get(url)
	if err != nil {
		return Pokemon{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, err
	}
	p.cache.Add(url, data)

	if err := json.Unmarshal(data, &pokemonData); err != nil {
		return Pokemon{}, err
	}

	return pokemonData, nil
}

func NewClient(cache *pokecache.Cache) PokeClient {
	var client PokeClient
	client.httpClient = &http.Client{}
	client.cache = cache
	client.BaseURL = "https://pokeapi.co/api/v2/"
	client.Config = APIConfig{}
	client.Pokedex = make(map[string]Pokemon)
	return client
}
