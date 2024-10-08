package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationResponse struct {
	Count    int        `json:"count"`
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
	Results  []Location `json:"results"`
}

func GetLocations(url string) (LocationResponse, error) {
	res, err := http.Get(url)
	if err != nil {
		return LocationResponse{}, fmt.Errorf("error getting locations: %w", err)
	}
	defer res.Body.Close()
	var locationRes LocationResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&locationRes); err != nil {
		return LocationResponse{}, fmt.Errorf("error reading response data: %w", err)
	}
	return locationRes, nil
}
