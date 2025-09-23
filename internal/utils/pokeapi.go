package pokeapi

import (
	"fmt"
	"io"
	"time"
	"encoding/json"
	"net/http"

	pokecache "github.com/louiehdev/pokedexcli/internal/cache"
)

type PokeClient struct {
	httpClient *http.Client
	cache      *pokecache.Cache
	BaseURL    string
	Config APIConfig
}

type APIConfig struct {
	Next     *string
	Previous *string
}

func NewClient(cache *pokecache.Cache) PokeClient {
	var client PokeClient
	client.httpClient = &http.Client{}
	client.cache = cache
	client.BaseURL = "https://pokeapi.co/api/v2/location-area"
	client.Config = APIConfig{}
	return client
}

type PokeLocationData struct {
	Next     *string
	Previous *string
	Results  []PokeLocationArea
}

type PokeLocationArea struct {
	Name string
	Url  string
}

func (p *PokeClient) GetPokeLocationData(url string) (PokeLocationData, error) {
	var locationData PokeLocationData

	if entry, exists := p.cache.PokeCache[url]; exists {
		json.Unmarshal(entry.Data, &locationData)
		return locationData, nil
	}

	res, err := p.httpClient.Get(url)
	if err != nil {
		return PokeLocationData{}, fmt.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return PokeLocationData{}, err
	}
	p.cache.PokeCache[url] = pokecache.CacheEntry{CreatedAt: time.Now(), Data: data}

	if err := json.Unmarshal(data, &locationData); err != nil {
		return PokeLocationData{}, err
	}
	/*
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&locationData); err != nil {
			return PokeLocationData{}, err
		}
	*/
	return locationData, nil
}
