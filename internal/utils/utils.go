package utils

import (
	"fmt"
	"net/http"
	"encoding/json"
)

type PokeLocationData struct {
	Next *string
	Previous *string
	Results []PokeLocationArea
}

type PokeLocationArea struct {
	Name string
	Url string
}

func GetPokeLocationData(url string) (PokeLocationData, error) {
	var locationData PokeLocationData
	res, err := http.Get(url)
	if err != nil {
		return locationData, fmt.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&locationData); err != nil {
		return PokeLocationData{}, err
	}

	return locationData, nil

}