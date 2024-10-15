package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tarikstupac/pokedex/internal/pokecache"
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

func GetLocations(url string, cache *pokecache.Cache) (LocationResponse, error) {
	var locationRes LocationResponse
	data, ok := cache.Get(url)
	if ok {
		if err := json.Unmarshal(data, &locationRes); err != nil {
			return LocationResponse{}, fmt.Errorf("error parsing data from cache: %w", err)
		}
		return locationRes, nil
	}
	res, err := http.Get(url)
	if err != nil {
		return LocationResponse{}, fmt.Errorf("error getting locations: %w", err)
	}
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return LocationResponse{}, fmt.Errorf("error reading data from response body: %w", err)
	}
	if err = cache.Add(url, data); err != nil {
		return LocationResponse{}, fmt.Errorf("error adding data to cache: %w", err)
	}

	if err = json.Unmarshal(data, &locationRes); err != nil {
		return LocationResponse{}, fmt.Errorf("error parsing data to json: %w", err)
	}
	return locationRes, nil
}
