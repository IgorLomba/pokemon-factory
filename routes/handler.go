package routes

import (
	"fmt"
	"net/http"

	"pokemon/controllers"
)

func RegisterRoutes(pokemon controllers.PokemonController) {
	http.HandleFunc(fmt.Sprintf("%s/generate", controllers.PokemonAPIPath), pokemon.GeneratePokemons())
	http.HandleFunc(fmt.Sprintf("%s/status", controllers.PokemonAPIPath), pokemon.CheckStatus())
	http.HandleFunc(fmt.Sprintf("%s/get", controllers.PokemonAPIPath), pokemon.GetPokemonsPolling())
	http.HandleFunc(fmt.Sprintf("%s/store", controllers.PokemonAPIPath), pokemon.LatestPokemon())
}
