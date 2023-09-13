package controllers

import (
	"net/http"

	"pokemon/models"
)

var pokemonHashMap = make(map[string]pokemonControllerHTTP)

type pokemonControllerHTTP struct {
	Pokemons      []models.Pokemon      `json:"pokemonsList"`
	RequestStatus *models.RequestStatus `json:"requestStatus"`
}

type PokemonController interface {
	GeneratePokemons() func(w http.ResponseWriter, r *http.Request)
	CheckStatus() func(w http.ResponseWriter, r *http.Request)
	GetPokemonsPolling() func(w http.ResponseWriter, r *http.Request)
	LatestPokemon() func(w http.ResponseWriter, r *http.Request)
}
