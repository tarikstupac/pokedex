package pokedexdata

import (
	"github.com/tarikstupac/pokedex/internal/pokeapi"
)

type Pokedex struct {
	pokemonCaught map[string]pokeapi.Pokemon
}

func NewPokedex() *Pokedex {
	p := Pokedex{pokemonCaught: map[string]pokeapi.Pokemon{}}
	return &p
}

func (p *Pokedex) Add(key string, val pokeapi.Pokemon) error {
	p.pokemonCaught[key] = val
	return nil
}
func (p *Pokedex) Get(key string) (pokeapi.Pokemon, bool) {
	entry, ok := p.pokemonCaught[key]
	if !ok {
		return pokeapi.Pokemon{}, ok
	}
	return entry, ok
}
