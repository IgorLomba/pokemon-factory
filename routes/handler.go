package routes

import (
	"fmt"
	"net/http"

	"pokemon/controllers"

	swaggerMiddleware "github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

const swaggerJSON = "swagger.json"

func RegisterRoutes(pokemon controllers.PokemonController) {
	swaggerHandler := swaggerHttpHandler()
	serveMux := http.NewServeMux()
	http.Handle(fmt.Sprintf("/%s", swaggerJSON), http.FileServer(http.Dir(".")))

	serveMux.Handle("/", swaggerHandler)
	serveMux.Handle("/doc", swaggerHandler)

	http.HandleFunc(fmt.Sprintf("%s/generate", controllers.PokemonAPIPath), pokemon.GeneratePokemons())
	http.HandleFunc(fmt.Sprintf("%s/status", controllers.PokemonAPIPath), pokemon.CheckStatus())
	http.HandleFunc(fmt.Sprintf("%s/get", controllers.PokemonAPIPath), pokemon.GetPokemonsPolling())
	http.HandleFunc(fmt.Sprintf("%s/store", controllers.PokemonAPIPath), pokemon.LatestPokemon())
	http.HandleFunc("/", serveMux.ServeHTTP)
	http.HandleFunc("/doc", serveMux.ServeHTTP)
}

func swaggerHttpHandler() http.Handler {
	router := mux.NewRouter()

	swaggerOptions := swaggerMiddleware.SwaggerUIOpts{SpecURL: swaggerJSON, Title: "Pokemon API"}
	ui := swaggerMiddleware.SwaggerUI(swaggerOptions, nil)
	router.Handle("/docs", ui)

	redocOptions := swaggerMiddleware.RedocOpts{SpecURL: swaggerJSON, Path: "doc"}
	doc := swaggerMiddleware.Redoc(redocOptions, nil)
	router.Handle("/doc", doc)

	return router
}
