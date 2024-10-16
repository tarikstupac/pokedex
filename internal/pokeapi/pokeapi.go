package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tarikstupac/pokedex/internal/pokecache"
)

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

func GetEncounters(param string, cache *pokecache.Cache) (LocationDetailResponse, error) {
	var locationDetailRes LocationDetailResponse
	url := "https://pokeapi.co/api/v2/location-area/" + param
	data, ok := cache.Get(url)
	if ok {
		if err := json.Unmarshal(data, &locationDetailRes); err != nil {
			return LocationDetailResponse{}, fmt.Errorf("error parsing data from cache: %w", err)
		}
		return locationDetailRes, nil
	}
	res, err := http.Get(url)
	if err != nil {
		return LocationDetailResponse{}, fmt.Errorf("error getting locations: %w", err)
	}
	if res.StatusCode > 299 {
		return LocationDetailResponse{}, fmt.Errorf("no data found for location %v", param)
	}
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return LocationDetailResponse{}, fmt.Errorf("error reading data from response body: %w", err)
	}
	if err = cache.Add(url, data); err != nil {
		return LocationDetailResponse{}, fmt.Errorf("error adding data to cache: %w", err)
	}

	if err = json.Unmarshal(data, &locationDetailRes); err != nil {
		return LocationDetailResponse{}, fmt.Errorf("error parsing data to json: %w", err)
	}
	return locationDetailRes, nil
}

func GetPokemon(param string, cache *pokecache.Cache) (Pokemon, error) {
	var pokemonRes Pokemon
	url := "https://pokeapi.co/api/v2/pokemon/" + param
	data, ok := cache.Get(url)
	if ok {
		if err := json.Unmarshal(data, &pokemonRes); err != nil {
			return Pokemon{}, fmt.Errorf("error parsing data from cache: %w", err)
		}
		return pokemonRes, nil
	}

	res, err := http.Get(url)

	if err != nil {
		return Pokemon{}, fmt.Errorf("error getting pokemon: %w", err)
	}

	if res.StatusCode > 299 {
		return Pokemon{}, fmt.Errorf("no data found for pokemon: %v", param)
	}
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, fmt.Errorf("error reading data from response body: %w", err)
	}
	if err = cache.Add(url, data); err != nil {
		return Pokemon{}, fmt.Errorf("error adding data to cache: %w", err)
	}
	if err = json.Unmarshal(data, &pokemonRes); err != nil {
		return Pokemon{}, fmt.Errorf("error parsing data to json: %w", err)
	}
	return pokemonRes, nil
}
